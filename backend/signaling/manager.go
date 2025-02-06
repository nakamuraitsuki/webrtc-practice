package signaling

import (
	"sync"

	"github.com/gorilla/websocket"
)

type SignalingManager struct {
	clients		map[*websocket.Conn]string	// クライアントの接続状況（Conn → ID）
	clientsByID map[string]*websocket.Conn	// クライアントのIDごとの接続情報
	offerId		string						// 処理中のoffer識別子
	sdpData		map[string]string			// SDPデータの保存
	candidateData map[string][]string		// ICE候補の保存
	mu 			sync.Mutex					// データ競合を防ぐためのミューテックス
}

func NewSignalingManager() *SignalingManager {
	return &SignalingManager{
		clients: 		make(map[*websocket.Conn]string),
		clientsByID: 	make(map[string]*websocket.Conn),
		sdpData: 		make(map[string]string),
		candidateData: 	make(map[string][]string),
	}
}

//新しいクライアントの追加
func (sm *SignalingManager) AddClient(conn *websocket.Conn, id string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.clients[conn] = id
	sm.clientsByID[id] = conn
}

//クライアントの削除
func (sm *SignalingManager) RemoveClient(conn *websocket.Conn) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	id, exists := sm.clients[conn]
	if exists {
		delete(sm.clients, conn)
		delete(sm.clientsByID, id)
	}
}

//SDP情報保存
func (sm *SignalingManager) SaveSDP(id, sdp string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.sdpData[id] = sdp
}

//ICE Candidate保存
func (sm *SignalingManager) SaveCandidate(id, candidate string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.candidateData[id] = append(sm.candidateData[id], candidate)
}

//Candidateを取得
func (sm *SignalingManager) GetCandidate(id string) ([]string, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	candidates, exists := sm.candidateData[id]

	return candidates, exists
}

//offerIdをセット
func (sm *SignalingManager) SetOfferID(id string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.offerId = id
}

//offerId取得
func (sm *SignalingManager) GetOfferID() string {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	return sm.offerId
}

//offerIdをリセット
func (sm *SignalingManager) ResetOfferID() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.offerId = ""
}

//offerIdの存在確認
func (sm *SignalingManager) IsOfferIDSet() bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	return sm.offerId != ""
}