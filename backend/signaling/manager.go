package signaling

import (
	"sync"

	"github.com/gorilla/websocket"
)

type SignalingManager struct {
	clients		map[*websocket.Conn]string	// クライアントの接続状況（Conn → ID）
	clientsByID map[string]*websocket.Conn	// クライアントのIDごとの接続情報
	broadcast	chan []byte					// ブロードキャスト用のチャネル
	offerId		string						// 処理中のoffer識別子
	sdpData		map[string]string			// SDPデータの保存
	candidateData map[string][]string		// ICE候補の保存
	mu 			sync.Mutex					// データ競合を防ぐためのミューテックス
}

func NweSignalingManager() *SignalingManager {
	return &SignalingManager{
		clients: 		make(map[*websocket.Conn]string),
		clientsByID: 	make(map[string]*websocket.Conn),
		broadcast: 		make(chan []byte),
		sdpData: 		make(map[string]string),
		candidateData: 	make(map[string][]string),
	}
}

//新しいクライアントの追加
func (wm *SignalingManager) AddClient(conn *websocket.Conn, id string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	wm.clients[conn] = id
	wm.clientsByID[id] = conn
}

//クライアントの削除
func (wm *SignalingManager) RemoveClient(conn *websocket.Conn) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	id, exists := wm.clients[conn]
	if exists {
		delete(wm.clients, conn)
		delete(wm.clientsByID, id)
	}
}

//SDP情報保存
func (wm *SignalingManager) SaveSDP(id, sdp string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	wm.sdpData[id] = sdp
}

//ICE Candidate保存
func (wm *SignalingManager) SaveCandidate(id, candidate string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	wm.candidateData[id] = append(wm.candidateData[id], candidate)
}

//Candidateを取得
func (wm *SignalingManager) GetCandidate(id string) ([]string, bool) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	candidates, exists := wm.candidateData[id]

	return candidates, exists
}

//offerIdをセット
func (wm *SignalingManager) SetOfferID(id string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	wm.offerId = id
}

//offerId取得
func (wm *SignalingManager) GetOfferID() string {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	return wm.offerId
}

//offerIdをリセット
func (wm *SignalingManager) ResetOfferID() {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	wm.offerId = ""
}

//offerIdの存在確認
func (wm *SignalingManager) IsOfferIDSet() bool {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	return wm.offerId != ""
}