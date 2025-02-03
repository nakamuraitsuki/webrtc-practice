package signaling

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
)

type WebSocketService struct {
	Port string
}

type WebSocketHandler struct {
	service *WebSocketService
}

func NewWebSocketHandler(service *WebSocketService) *WebSocketHandler {
	return &WebSocketHandler{service: service}
}

// WebSocket接続の処理
func (h *WebSocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	// WebSocketの接続を確立
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	// 接続時の処理
	log.Println("New WebSocket connection established on port", h.service.Port)
}