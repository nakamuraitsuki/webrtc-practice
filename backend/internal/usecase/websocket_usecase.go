package usecase

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"example.com/webrtc-practice/internal/domain/service"
	"example.com/webrtc-practice/internal/infrastructure/repository_impl"
	offerservice "example.com/webrtc-practice/internal/infrastructure/service_impl/offer_service"
	websocketbroadcast "example.com/webrtc-practice/internal/infrastructure/service_impl/websocket_broadcast"
	websocketmanager "example.com/webrtc-practice/internal/infrastructure/service_impl/websocket_manager"
	"github.com/gorilla/websocket"
)

type IWebsocketUsecase struct {
	repo repository_impl.WebsocketRepositoryImpl
	wm   service.WebsocketManager
	br   service.WebSocketBroadcastService
	o    service.OfferService
}

func NewWebsocketUsecase() IWebsocketUsecase {
	return IWebsocketUsecase{
		repo: *repository_impl.NewWebsocketRepositoryImpl(),
		wm:   websocketmanager.NewWebsocketManager(),
		br:   websocketbroadcast.NewBroadcast(),
		o:    offerservice.NewOfferService(),
	}
}

// RegisterClientは新しいクライアントを登録（repoのラップ）
func (u *IWebsocketUsecase) RegisterClient(conn service.WebSocketConnection) error {
	return u.wm.RegisterConnection(conn)
}

// メッセージ受信待機（ユーザー　->　ブロードキャスト）
func (u *IWebsocketUsecase) ListenForMessages(conn service.WebSocketConnection) {
	// ID初期化
	var clientID string
	clientID = ""

	// メッセージ受信ループ
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			u.wm.DeleteConnection(conn)
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

			if u.wm.ExistsByID(id) {
				// 既に登録されている場合は、今つなごうとしているコネクションを削除
				u.wm.DeleteConnection(conn)
				log.Println("Client with ID already exists. Connection closed.")
				break
			}

			u.wm.RegisterID(conn, id)
		}
		u.br.Send(message)
	}
}

// メッセージ待ち（ブロードキャスト　->　サーバー）
func (u *IWebsocketUsecase) ProcessMessage() {
	for {
		message := u.br.Receive()
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
	client, err := u.wm.GetConnectionByID(id)
	if err != nil {
		log.Println("Client not found:", err)
		return
	}

	// もしofferしている人がいなかったら
	if !u.o.IsOffer() {
		// 現在offer中のIDを更新
		u.o.SetOffer(id)
		// offerをコールバック（送り主がofferを送ることを期待する）
		resultData["type"] = "offer"
		bytes := u.jsonToBytes(resultData)
		// 送信
		u.sendMessage(client, bytes)
		return
	} else if u.o.IsOfferID(id) { // offer中なのが自分だったら
		// 重複なので何もしない
		return
	}

	// もし自分以外のofferしている人がいたら。

	// anser待機中の人が送ったofferを整形（offerを受け取った相手がanswerを送ることを期待する）
	resultData["type"] = "offer"
	resultData["sdp"], err = u.repo.GetSDPByID(u.o.GetOffer())
	if err != nil {
		log.Println("SDP not found:", err)
		return
	}
	resultData["target_id"] = u.o.GetOffer()
	bytes := u.jsonToBytes(resultData)

	// 送信
	u.sendMessage(client, bytes)
}

func (u *IWebsocketUsecase) offer(data map[string]any) {
	// offerの送り主のSDPを保存
	fmt.Println("[Offer]")
	id := data["id"].(string)
	sdp, _ := json.Marshal(data["sdp"])
	u.repo.SaveSDP(id, string(sdp))
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

	client, err := u.wm.GetConnectionByID(target_id)
	if err != nil {
		log.Println("Client not found:", err)
		return
	}

	bytes := u.jsonToBytes(resultData)
	u.sendMessage(client, bytes)
}

func (u *IWebsocketUsecase) sendCandidate(data map[string]any) {
	returnData := make(map[string]string)
	id := u.o.GetOffer()

	if !u.repo.ExistsCandidateByID(id) {
		return
	}

	answerId := data["id"].(string)
	// クライアントの取得（repo）
	client, err := u.wm.GetConnectionByID(answerId)
	if err == nil {
		log.Println("Client not found:", answerId)
		return
	}

	fmt.Println("candidate受け取り")
	fmt.Println("[Candidate]")
	returnData["type"] = "candidate"

	candidate, err := u.repo.GetCandidatesByID(id)
	if err != nil {
		log.Println("Candidate not found:", err)
		return
	}
	returnData["candidate"] = strings.Join(candidate, "|")

	bytes := u.jsonToBytes(returnData)

	// 送信
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
		if client, err := u.wm.GetConnectionByID(target_id); err == nil {
			// 相手が接続中
			fmt.Println("[Candidate]")
			resultData["type"] = "candidate"
			resultData["candidate"] = candidate
			bytes := u.jsonToBytes(resultData)

			// 送信
			u.sendMessage(client, bytes)
			return
		}
	}

	// 相手が還沒來 -> 保存
	if !u.repo.ExistsCandidateByID(id) {
		err := u.repo.SaveCandidate(id, candidate)
		if err != nil {
			log.Println("Error saving candidate:", err)
			return
		}
	} else {
		err := u.repo.AddCandidate(id, candidate)
		if err != nil {
			log.Println("Error adding candidate:", err)
			return
		}
	}
}

// message送信
func (u *IWebsocketUsecase) sendMessage(client service.WebSocketConnection, bytes []byte) {
	err := client.WriteMessage(websocket.TextMessage, bytes)
	if err != nil {
		log.Println(err)
		u.wm.DeleteConnection(client)
		// ハンドラ内で defer conn.Close() の使用を期待してコネクションの閉鎖はしない
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
