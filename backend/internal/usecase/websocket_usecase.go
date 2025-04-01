package usecase

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/domain/repository"
	"example.com/webrtc-practice/internal/domain/service"
	"example.com/webrtc-practice/internal/infrastructure/repository_impl"
	offerservice "example.com/webrtc-practice/internal/infrastructure/service_impl/offer_service"
	websocketbroadcast "example.com/webrtc-practice/internal/infrastructure/service_impl/websocket_broadcast"
	websocketmanager "example.com/webrtc-practice/internal/infrastructure/service_impl/websocket_manager"
	"github.com/gorilla/websocket"
)

type IWebsocketUsecase struct {
	repo repository.IWebsocketRepository
	wm   service.WebsocketManager
	br   service.WebSocketBroadcastService
	o    service.OfferService
}

func NewWebsocketUsecase() IWebsocketUsecase {
	return IWebsocketUsecase{
		repo: repository_impl.NewWebsocketRepositoryImpl(),
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
		// メッセージ抽出
		// 文字列 -> json
		var jsonStr = string(message)

		// json -> map[string]any
		var data entity.Message
		err = json.Unmarshal([]byte(jsonStr), &data)
		if err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		// 初回ID登録
		if clientID == "" {


			// idの取得
			id := data.ID
			clientID = data.ID

			if u.wm.ExistsByID(id) {
				// 既に登録されている場合は、今つなごうとしているコネクションを削除
				u.wm.DeleteConnection(conn)
				log.Println("Client with ID already exists. Connection closed.")
				break
			}

			u.wm.RegisterID(conn, id)
			u.repo.CreateClient(id)
		}
		u.br.Send(data)
	}
	u.o.ClearOffer()
}

// メッセージ待ち（ブロードキャスト　->　サーバー）
func (u *IWebsocketUsecase) ProcessMessage() {
	for {
		message := u.br.Receive()

		// 処理の分岐
		msgType := message.Type

		switch msgType {
		case "connect":
			u.connect(message)
		case "offer":
			u.offer(message)
		case "answer":
			u.answer(message)
		case "candidateAdd":
			u.candidateAdd(message)
		}
	}
}

func (u *IWebsocketUsecase) connect(message entity.Message) {
	resultData := make(map[string]string)

	// メッセージの送り主を取得
	id := message.ID
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

func (u *IWebsocketUsecase) offer(message entity.Message) {
	// offerの送り主のSDPを保存
	fmt.Println("[Offer]")
	id := message.ID
	sdp, _ := json.Marshal(message.SDP)
	u.repo.SaveSDP(id, string(sdp))
}

func (u *IWebsocketUsecase) answer(message entity.Message) {
	// offerの送り主にanswerを返す
	u.sendAnswer(message)

	// answerの送り主にcandidateを送る
	u.sendCandidate(message)
}

func (u *IWebsocketUsecase) sendAnswer(message entity.Message) {
	fmt.Println("[Answer]")
	resultData := make(map[string]string)
	resultData["type"] = "answer"
	target_id := message.TargetID
	sdp, _ := json.Marshal(message.SDP)
	resultData["sdp"] = string(sdp)

	client, err := u.wm.GetConnectionByID(target_id)
	if err != nil {
		log.Println("Client not found:", err)
		return
	}

	bytes := u.jsonToBytes(resultData)
	u.sendMessage(client, bytes)
}

func (u *IWebsocketUsecase) sendCandidate(message entity.Message) {
	returnData := make(map[string]string)
	id := u.o.GetOffer()

	if !u.repo.ExistsCandidateByID(id) {
		return
	}

	answerId := message.ID
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

func (u *IWebsocketUsecase) candidateAdd(message entity.Message) {
	fmt.Println("[Candidate Add]")
	resultData := make(map[string]string)

	// 相手が通話中なら、candidateDataに入れずに直接送る
	id := message.ID
	candidateByte, _ := json.Marshal(message.Candidate)
	candidate := string(candidateByte)

	target_id := message.TargetID
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
