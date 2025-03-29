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
	broadcast            = make(chan []byte) // []byteだが、Message型を作りたい気持ちがある
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

// RegisterClientは新しいクライアントを登録
func (u *IWebsocketUsecase) RegisterClient(conn *websocket.Conn) error {
	// connectionの存在を登録(repo)
	// ミューテーションロックを使用して、同時アクセスを防止
	u.mu.Lock()
	defer u.mu.Unlock()

	// 重複登録を避ける
	if _, exists := clients[conn]; exists {
		return errors.New("client already registered")
	}

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

			// idの取得
			id := data["id"].(string)
			clientID = id

			// 重複登録の確認
			u.mu.Lock()
			if _, exists := clientsByID[id]; exists {
				u.mu.Unlock()
				log.Println("Error: User ID already registered")
				conn.Close()
				return
			}

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
			u.connect(data)
		case "offer":
			u.offer(data)
		case "answer":
			u.answer(data)
		case "candidateAdd":
			u.candidateAdd(data)
		}
	}
}

func (u *IWebsocketUsecase) connect(data map[string]any) {
	resultData := make(map[string]string)

	// メッセージの送り主を取得
	id := data["id"].(string)
	// IDからクライアントを取得(repo)
	client := clientsByID[id]

	// もしofferしている人がいなかったら
	if len(offerId) == 0 {
		// 現在offer中のIDを更新
		offerId = id
		// offerをコールバック（送り主がofferを送ることを期待する）
		resultData["type"] = "offer"
		bytes := u.jsonToBytes(resultData)
		// 送信（conn依存）
		u.sendMessage(client, bytes)
		return
	} else if id == offerId {// offer中なのが自分だったら
		// 重複なので何もしない
		return
	}

	// もし自分以外のofferしている人がいたら。

	// anser待機中の人が送ったofferを整形（offerを受け取った相手がanswerを送ることを期待する）
	resultData["type"] = "offer"
	resultData["sdp"] = sdpData[offerId]
	resultData["target_id"] = offerId
	bytes := u.jsonToBytes(resultData)

	// 送信（conn依存）
	u.sendMessage(client, bytes)
}

func (u *IWebsocketUsecase) offer(data map[string]any) {
	// offerの送り主のSDPを保存
	fmt.Println("[Offer]")
	id := data["id"].(string)
	sdp, _ := json.Marshal(data["sdp"])
	// 受け取ったSDPを保存(repo)
	sdpData[id] = string(sdp)
}

func (u *IWebsocketUsecase) answer(data map[string]any) {
	// offerの送り主にanswerを返す
	u.sendAnswer(data)

	// answerの送り主にcandidateを送る
	u.sendCandidate(data)
}

func (u *IWebsocketUsecase) sendAnswer(data map[string]any) {
	fmt.Println("[Answer]")
	resultData := make(map[string]string)
	resultData["type"] = "answer"
	target_id := data["target_id"].(string)
	sdp, _ := json.Marshal(data["sdp"])
	resultData["sdp"] = string(sdp)

	client := clientsByID[target_id]
	bytes := u.jsonToBytes(resultData)
	u.sendMessage(client, bytes)
}

func (u *IWebsocketUsecase) sendCandidate(data map[string]any) {
	returnData := make(map[string]string)
	id := offerId
	// 保存されているcandidateの有無を確認（repo）
	if _, ok := candidateData[id]; !ok {
		return
	}

	answerId := data["id"].(string)
	// クライアントの取得（repo）
	client := clientsByID[answerId]
	fmt.Println("candidate受け取り")
	fmt.Println("[Candidate]")
	returnData["type"] = "candidate"
	returnData["candidate"] = strings.Join(candidateData[id], "|")
	bytes := u.jsonToBytes(returnData)

	// 送信（conn依存）
	u.sendMessage(client, bytes)

}

func (u *IWebsocketUsecase) candidateAdd(data map[string]any) {
	fmt.Println("[Candidate Add]")
	resultData := make(map[string]string)

	// 相手が通話中なら、candidateDataに入れずに直接送る
	id := data["id"].(string)
	candidateByte, _ := json.Marshal(data["candidate"])
	candidate := string(candidateByte)

	target_id := data["target_id"].(string)
	if target_id != "" {
		if client, ok2 := clientsByID[target_id]; ok2 {
			// 相手が接続中
			fmt.Println("[Candidate]")
			resultData["type"] = "candidate"
			resultData["candidate"] = candidate
			bytes := u.jsonToBytes(resultData)

			// 送信（conn依存）
			u.sendMessage(client, bytes)
			return
		}
	}

	// 相手が還沒來 -> 保存
	// candidateの存在確認（repo）
	if _, ok := candidateData[id]; !ok {
		// candidateの保存(repo)
		candidateData[id] = []string{candidate}
	} else {
		// candidateの追加(repo)
		candidateData[id] = append(candidateData[id], candidate)
	}
}

// message送信（conn 依存）
func (u *IWebsocketUsecase) sendMessage(client *websocket.Conn, bytes []byte) {
	err := client.WriteMessage(websocket.TextMessage, bytes)
	if err != nil {
		log.Println(err)
		client.Close()
		delete(clients, client)
	}
}

func (u *IWebsocketUsecase) jsonToBytes(result map[string]string) []byte {
	jsonText, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	bytes := []byte(jsonText)
	return bytes
}
