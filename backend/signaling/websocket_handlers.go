package signaling

import (
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
	manager *SignalingManager
	upgrader websocket.Upgrader
}

func NewWebSocketHandler(service *WebSocketService, sm *SignalingManager, upgrader websocket.Upgrader) *WebSocketHandler {
	return &WebSocketHandler{
		service: service,
		manager: sm,
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
	
	h.manager.ResetOfferID()
	return nil
}