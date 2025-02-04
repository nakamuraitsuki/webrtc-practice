package signaling

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WebSocketService struct {
	Port string
}

type WebSocketHandler struct {
	service *WebSocketService
	manager *WebSocketManager
	upgrader websocket.Upgrader
}

func NewWebSocketHandler(service *WebSocketService, upgrader websocket.Upgrader) *WebSocketHandler {
	return &WebSocketHandler{
		service: service,
		upgrader: upgrader,
	}
}

func (h *WebSocketHandler) Register(g *echo.Group) {
	g.GET("", h.HandleWebSocket)
}

// webSocket接続の処理
func (h *WebSocketHandler) HandleWebSocket(c echo.Context) error{
	req := c.Request()
	resp := c.Response().Writer
	
	conn, err := h.upgrader.Upgrade(resp, req, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)

		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to establish WebSocket connection",
		})
	}
	defer conn.Close()

	// 接続時の処理
	log.Println("New WebSocket connection established on port", h.service.Port)

	//メッセージ受信待機ループ
	for {
		//クライアントからメッセージを受け取る
		_, message, err := conn.ReadMessage()

		//新規かどうか確認
		if _, ok := h.manager.clients[conn]; !ok {
			//新規接続
			var jsonStr = string(message)
			var data map[string]interface{}
			err := json.Unmarshal([]byte(jsonStr), &data)
			if err != nil {
				panic(err)
			}

			//id登録
			id := data["id"].(string)
			h.manager.AddClient(conn, id)
		}

		if err != nil {
			log.Println(err)
			h.manager.RemoveClient(conn)
			break
		}

		//ブロードキャストにpush
		h.manager.broadcast <- message
	}
	
	h.manager.ResetOfferID()
	return nil
}