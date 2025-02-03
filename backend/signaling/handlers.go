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
}

func NewWebSocketHandler(service *WebSocketService) *WebSocketHandler {
	return &WebSocketHandler{service: service}
}

func (h *WebSocketHandler) Register(g *echo.Group) {
	g.GET("", h.HandleConnection)
}

// webSocket接続の処理
func (h *WebSocketHandler) HandleConnection(c echo.Context) error{
	// WebSocketの接続を確立
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	req := c.Request()

	resp := c.Response().Writer
	
	conn, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)

		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to establish WebSocket connection",
		})
	}
	defer conn.Close()

	// 接続時の処理
	log.Println("New WebSocket connection established on port", h.service.Port)

	return nil
}