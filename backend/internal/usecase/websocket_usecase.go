package usecase

import (
	"fmt"
	"log"

	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/domain/repository"
	"example.com/webrtc-practice/internal/domain/service"
)

type IWebsocketUsecaseInterface interface {
	RegisterClient(conn service.WebSocketConnection) error
	ListenForMessages(conn service.WebSocketConnection)
	ProcessMessage()
	Connect(message entity.Message)
	Offer(message entity.Message)
	Answer(message entity.Message)
	Candidate(message entity.Message)
	CandidateAdd(message entity.Message) bool
	SendCandidate(message entity.Message)
}

type IWebsocketUsecase struct {
	repo repository.IWebsocketRepository
	wm   service.WebsocketManager
	br   service.WebSocketBroadcastService
	o    service.OfferService
}

func NewWebsocketUsecase(
	repo repository.IWebsocketRepository,
	wm service.WebsocketManager,
	br service.WebSocketBroadcastService,
	o service.OfferService,
) *IWebsocketUsecase {
	return &IWebsocketUsecase{
		repo: repo,
		wm:   wm,
		br:   br,
		o:    o,
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
			u.repo.DeleteClient(clientID)
			break
		}


		// 初回ID登録
		if clientID == "" {
			// idの取得
			id := message.ID
			clientID = message.ID

			if u.wm.ExistsByID(id) {
				// 既に登録されている場合は、今つなごうとしているコネクションを削除
				u.wm.DeleteConnection(conn)
				log.Println("Client with ID already exists. Connection closed.")
				break
			}

			err = u.wm.RegisterID(conn, id)
			if err != nil {
				log.Println("Failed to register ID:", err)
				break
			}
			u.repo.CreateClient(id)
		}
		u.br.Send(message)
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
			u.Connect(message)
		case "offer":
			u.Offer(message)
		case "answer":
			u.Answer(message)
		case "candidateAdd":
			u.Candidate(message)
		default:
			log.Println("Unknown message type:", msgType)
		}
	}
}

func (u *IWebsocketUsecase) Connect(message entity.Message) {
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
		resultData := entity.Message{
			ID:   id,
			Type: "offer",
		}
		// 送信
		client.WriteMessage(resultData)
		return
	} else if u.o.IsOfferID(id) { // offer中なのが自分だったら
		// 重複なので何もしない
		return
	}

	// もし自分以外のofferしている人がいたら。

	// anser待機中の人が送ったofferを整形（offerを受け取った相手がanswerを送ることを期待する）
	sdp, err := u.repo.GetSDPByID(u.o.GetOffer())
	if err != nil {
		log.Println("SDP not found:", err)
		return
	}
	targetID := u.o.GetOffer()

	resultData := entity.Message{
		ID:       id,
		Type:     "offer",
		SDP:      sdp,
		TargetID: targetID,
	}
	// 送信
	client.WriteMessage(resultData)
}

func (u *IWebsocketUsecase) Offer(message entity.Message) {
	// offerの送り主のSDPを保存
	fmt.Println("[Offer]")
	u.repo.SaveSDP(message.ID, message.SDP)
}

func (u *IWebsocketUsecase) Answer(message entity.Message) {
	fmt.Println("[Answer]")
	targetID := message.TargetID

	client, err := u.wm.GetConnectionByID(targetID)
	if err != nil {
		log.Println("Client not found:", err)
		return
	}

	resultData := entity.Message{
		ID : message.ID,
		Type: "answer",
		SDP: message.SDP,
		TargetID: targetID,
	}

	client.WriteMessage(resultData)
}

func (u *IWebsocketUsecase) Candidate(message entity.Message) {
	fmt.Println("[Candidate]")
	// candidateを保存
	result := u.CandidateAdd(message); 
	if result {
		u.SendCandidate(message)
	}
}

// 別途送信の必要がある場合のみtrueを返し、sendCandidateを呼ぶ
func (u *IWebsocketUsecase) CandidateAdd(message entity.Message) bool {
	fmt.Println("[Candidate Add]")
	resultData := entity.Message{}
	
	id := message.ID
	candidate := message.Candidate
	targetID := message.TargetID
	
	// 通信相手の認知が完了している -> 直接送信
	if targetID != "" {
		if client, err := u.wm.GetConnectionByID(targetID); err == nil {
			fmt.Println("[Candidate]")
			resultData.Type = "candidate"
			resultData.ID = id
			resultData.Candidate = candidate
			client.WriteMessage(resultData)
		}
	}
	
	// 保存（初回 or 追加）
	if !u.repo.ExistsCandidateByID(id) {
		if err := u.repo.SaveCandidate(id, candidate); err != nil {
			log.Println("Error saving candidate:", err)
			return false
		}
	} else {
		if err := u.repo.AddCandidate(id, candidate); err != nil {
			log.Println("Error adding candidate:", err)
			return false
		}
	}
	
	// Answerer から送られてきた candidateAdd であれば、
	// Offerer 側の candidate を送る
	if u.o.IsOfferID(message.TargetID) {
		// Offerer から answerer に candidate を送る
		return true
	}

	return false
}
	
func (u *IWebsocketUsecase) SendCandidate(message entity.Message) {
	returnData := entity.Message{}
	// 送信元の名義
	id := u.o.GetOffer()

	if !u.repo.ExistsCandidateByID(id) {
		return
	}

	answerId := message.ID
	// クライアントの取得（repo）
	client, err := u.wm.GetConnectionByID(answerId)
	if err != nil {
		log.Println("Client not found:", answerId)
		return
	}	

	fmt.Println("candidate受け取り")
	fmt.Println("[Candidate]")
	returnData.Type = "candidate"

	candidate, err := u.repo.GetCandidatesByID(id)
	if err != nil {
		log.Println("Candidate not found:", err)
		return
	}
	returnData.Candidate = candidate

	// 送信
	client.WriteMessage(returnData)
}
