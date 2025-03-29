package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	clients              = make(map[*websocket.Conn]string)
	clientsByID          = make(map[string]*websocket.Conn)
	broadcast            = make(chan []byte)
	sdpData              = make(map[string]string)
	candidateData        = make(map[string][]string)
	offerId       string = ""
)

type IWebsocketUsecase struct {
	mu *sync.Mutex
}

func NewWebsocketUsecase(
	mu *sync.Mutex,
) IWebsocketUsecase {
	return IWebsocketUsecase{
		mu: mu,
	}
}

// RegisterClientは新しいクライアントを登録し、メッセージ受信のゴルーチンを開始
func (u *IWebsocketUsecase) RegisterClient(conn *websocket.Conn) error {
	// ミューテーションロックを使用して、同時アクセスを防止
	u.mu.Lock()
	defer u.mu.Unlock()

	// 重複登録を避ける
	if _, exists := clients[conn]; exists {
		return errors.New("client already registered")
	}

	// connectionの存在を登録(repo)
	clients[conn] = ""

	return nil
}

// メッセージ受信待機（ユーザー　->　ブロードキャスト）
func (u *IWebsocketUsecase) ListenForMessages(conn *websocket.Conn) {
	// ID初期化
	var clientID string
	clientID = ""

	// メッセージ受信ループ
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			u.HandleClientDisconnection(conn)
			break
		}

		// 初回ID登録
		if clientID == "" {
			// メッセージ抽出
			// 文字列 -> json
			var jsonStr = string(message)

			// json -> map[string]any
			var data map[string]any
			err := json.Unmarshal([]byte(jsonStr), &data)
			if err != nil {
				log.Println("Error unmarshalling message:", err)
				continue
			}

			id := data["id"].(string)
			clientID = id

			// 新規ユーザー本登録(repo)
			clients[conn] = id
			clientsByID[id] = conn
		}

		broadcast <- message
	}
}

// クライアント切断時の処理
func (u *IWebsocketUsecase) HandleClientDisconnection(conn *websocket.Conn) {
	// ミューテーションロックを使用して、同時アクセスを防止
	u.mu.Lock()
	defer u.mu.Unlock()

	// 削除処理(repo)
	if id, ok := clients[conn]; ok {
		delete(clientsByID, id)
		delete(clients, conn)
	}

	// ハンドラ内での defer conn.Close() の使用を期待して、コネクション閉鎖はしない
}

// メッセージ待ち（ブロードキャスト　->　サーバー）
func (u *IWebsocketUsecase) ProcessMessage() {
	for {
		message := <-broadcast
		// text -> json
		var jsonStr = string(message)
		fmt.Println(jsonStr)
		var data map[string]any
		err := json.Unmarshal([]byte(jsonStr), &data)
		if err != nil {
			panic(err)
		}

		// 処理の分岐
		msgDataType := data["type"].(string)

		switch msgDataType {
		case "connect":
			connect(data)
		case "offer":
			offer(data)
		case "answer":
			answer(data)
		case "candidateAdd":
			candidateAdd(data)
		}
	}
}

func connect(data map[string]any) {
	resultData := make(map[string]string)

	// offerの送り主を取得
	id := data["id"].(string)
	client := clientsByID[id]

	// offerを送ってもらう
	if len(offerId) == 0 {
		offerId = id
		resultData["type"] = "offer"
		bytes := jsonToBytes(resultData)
		sendMessage(client, bytes)
		return
	} else if id == offerId {
		// 重複
		return
	}

	// offerを送る
	resultData["type"] = "offer"
	resultData["sdp"] = sdpData[offerId]
	resultData["target_id"] = offerId
	bytes := jsonToBytes(resultData)
	sendMessage(client, bytes)
}

func offer(data map[string]any) {
	fmt.Println("[Offer]")
	id := data["id"].(string)
	sdp, _ := json.Marshal(data["sdp"])
	sdpData[id] = string(sdp)
}

func answer(data map[string]any) {
	// offerの送り主にanswerを返す
	sendAnswer(data)

	// answerの送り主にcandidateを送る
	sendCandidate(data)
}

func sendAnswer(data map[string]any) {
	fmt.Println("[Answer]")
	resultData := make(map[string]string)
	resultData["type"] = "answer"
	target_id := data["target_id"].(string)
	sdp, _ := json.Marshal(data["sdp"])
	resultData["sdp"] = string(sdp)

	client := clientsByID[target_id]
	bytes := jsonToBytes(resultData)
	sendMessage(client, bytes)
}

func sendCandidate(data map[string]any) {
	returnData := make(map[string]string)
	id := offerId
	if _, ok := candidateData[id]; !ok {
		return
	}

	answerId := data["id"].(string)
	client := clientsByID[answerId]
	fmt.Println("candidate受け取り")
	fmt.Println("[Candidate]")
	returnData["type"] = "candidate"
	returnData["candidate"] = strings.Join(candidateData[id], "|")
	bytes := jsonToBytes(returnData)
	sendMessage(client, bytes)

}

func candidateAdd(data map[string]any) {
	fmt.Println("[Candidate Add]")
	resultData := make(map[string]string)

	// 相手が已經接続的話、candidateDataに入れずに直接送る
	id := data["id"].(string)
	candidateByte, _ := json.Marshal(data["candidate"])
	candidate := string(candidateByte)

	target_id := data["target_id"].(string)
	if target_id != "" {
		if client, ok2 := clientsByID[target_id]; ok2 {
			// 相手が有的話
			fmt.Println("[Candidate]")
			resultData["type"] = "candidate"
			resultData["candidate"] = candidate
			bytes := jsonToBytes(resultData)
			sendMessage(client, bytes)
			return
		}
	}

	// 相手が還沒來 -> 保存
	if _, ok := candidateData[id]; !ok {
		candidateData[id] = []string{candidate}
	} else {
		candidateData[id] = append(candidateData[id], candidate)
	}
}

// 訊息送信
func sendMessage(client *websocket.Conn, bytes []byte) {
	err := client.WriteMessage(websocket.TextMessage, bytes)
	if err != nil {
		log.Println(err)
		client.Close()
		delete(clients, client)
	}
}

func jsonToBytes(result map[string]string) []byte {
	jsonText, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	bytes := []byte(jsonText)
	return bytes
}
