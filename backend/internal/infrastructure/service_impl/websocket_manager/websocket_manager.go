package websocketmanager

import (
	"errors"
	"sync"

	"example.com/webrtc-practice/internal/domain/service"
)

type WebsocketManagerImpl struct {
	clients     map[service.WebSocketConnection]string
	clientsByID map[string]service.WebSocketConnection
	mu          *sync.Mutex
}

func NewWebsocketManager() service.WebsocketManager {
	return &WebsocketManagerImpl{
		clients:     make(map[service.WebSocketConnection]string),
		clientsByID: make(map[string]service.WebSocketConnection),
		mu:          &sync.Mutex{},
	}
}

func (wm *WebsocketManagerImpl) RegisterConnection(conn service.WebSocketConnection) error {
	// ミューテーションロックを使用して、同時アクセスを防止
	wm.mu.Lock()
	defer wm.mu.Unlock()

	// 重複登録を避ける
	if _, exists := wm.clients[conn]; exists {
		return errors.New("client already registered")
	}

	wm.clients[conn] = ""

	return nil
}

func (wm *WebsocketManagerImpl) RegisterID(conn service.WebSocketConnection, id string) {
	// ミューテーションロックを使用して、同時アクセスを防止
	wm.mu.Lock()
	defer wm.mu.Unlock()

	// ID登録
	wm.clients[conn] = id
	wm.clientsByID[id] = conn
}

func (wm *WebsocketManagerImpl) DeleteConnection(conn service.WebSocketConnection) error {
	// ミューテーションロックを使用して、同時アクセスを防止
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if id, exists := wm.clients[conn]; !exists {
		// コネクションが見つからない場合
		return errors.New("client not found")
	} else if id == "" {
		// IDが登録されていない場合
		delete(wm.clients, conn)
	} else {
		// IDが登録されている場合
		delete(wm.clients, conn)
		delete(wm.clientsByID, id)
	}

	// ハンドラ内での defer conn.Close() の使用を期待してコネクションの閉鎖はしない
	return nil
}

func (wm *WebsocketManagerImpl) ExistsByID(id string) bool {
	// ミューテーションロックを使用して、同時アクセスを防止
	wm.mu.Lock()
	defer wm.mu.Unlock()

	_, exists := wm.clientsByID[id]
	return exists
}

func (wm *WebsocketManagerImpl) GetConnectionByID(id string) (service.WebSocketConnection, error) {
	// ミューテーションロックを使用して、同時アクセスを防止
	wm.mu.Lock()
	defer wm.mu.Unlock()

	conn, exists := wm.clientsByID[id]
	if !exists {
		return nil, errors.New("client not found")
	}
	return conn, nil
}
