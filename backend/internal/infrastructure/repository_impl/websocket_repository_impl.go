package repository_impl

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

type WebsocketRepositoryImpl struct {
	clients       map[*websocket.Conn]string
	clientsByID   map[string]*websocket.Conn
	sdpData       map[string]string
	candidateData map[string][]string
	mu            sync.Mutex
}

func NewWebsocketRepositoryImpl() *WebsocketRepositoryImpl {
	return &WebsocketRepositoryImpl{
		clients:       make(map[*websocket.Conn]string),
		clientsByID:   make(map[string]*websocket.Conn),
		sdpData:       make(map[string]string),
		candidateData: make(map[string][]string),
		mu:            sync.Mutex{},
	}
}

func (wr *WebsocketRepositoryImpl) RegisterConnection(conn *websocket.Conn) error {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	// 重複登録を避ける
	if _, exists := wr.clients[conn]; exists {
		return errors.New("client already registered")
	}

	wr.clients[conn] = ""

	return nil
}

func (wr *WebsocketRepositoryImpl) RegisterID(conn *websocket.Conn, id string) {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	// ID登録
	wr.clients[conn] = id
	wr.clientsByID[id] = conn
}

func (wr *WebsocketRepositoryImpl) DeleteConnection(conn *websocket.Conn) error {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	if id, exists := wr.clients[conn]; !exists {
		// コネクションが見つからない場合
		return errors.New("client not found")
	} else if id == "" {
		// IDが登録されていない場合
		delete(wr.clients, conn)
	} else {
		// IDが登録されている場合
		delete(wr.clients, conn)
		delete(wr.clientsByID, id)
	}

	// ハンドラ内での defer conn.Close() の使用を期待してコネクションの閉鎖はしない
	return nil
}

func (wr *WebsocketRepositoryImpl) ExistsByID(id string) bool {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	_, exists := wr.clientsByID[id]
	return exists
}

func (wr *WebsocketRepositoryImpl) GetClientByID(id string) (*websocket.Conn, error) {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	conn, exists := wr.clientsByID[id]
	if !exists {
		return nil, errors.New("client not found")
	}
	return conn, nil
}

func (wr *WebsocketRepositoryImpl) SaveSDP(id string, sdp string) error {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	wr.sdpData[id] = sdp
	return nil
}

func (wr *WebsocketRepositoryImpl) GetSDPByID(id string) (string, error) {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	sdp, exists := wr.sdpData[id]
	if !exists {
		return "", errors.New("SDP not found")
	}
	return sdp, nil
}

func (wr *WebsocketRepositoryImpl) SaveCandidate(id string, candidate string) error {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	wr.candidateData[id] = []string{candidate}
	return nil
}

func (wr *WebsocketRepositoryImpl) AddCandidate(id string, candidate string) error {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	if _, exists := wr.candidateData[id]; !exists {
		return errors.New("candidates not found")
	}

	wr.candidateData[id] = append(wr.candidateData[id], candidate)
	return nil
}

func (wr *WebsocketRepositoryImpl) GetCandidatesByID(id string) ([]string, error) {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	candidates, exists := wr.candidateData[id]
	if !exists {
		return nil, errors.New("candidates not found")
	}
	return candidates, nil
}

func (wr *WebsocketRepositoryImpl) ExistsCandidateByID(id string) bool {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	_, exists := wr.candidateData[id]
	return exists
}