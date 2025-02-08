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
	userID := c.Get("userID").(string)
	//TODO: ユーザー認証機能と、ミドルウェアでID設定
	
	//接続を取得
	conn, err := h.upgrader.Upgrade(resp, req, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)

		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to establish WebSocket connection",
		})
	}
	defer conn.Close()

	err = h.manager.AddClient(conn, userID);
	if err != nil {
		// エラーが発生した場合、クライアントにエラーメッセージを返す
		log.Println("Error adding client:", err)

		// IDがすでに登録されている場合は409 Conflict
		if err.Error() == "client already exists with the given ID" {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "Client with this ID already exists",
			})
		}
		
		// その他のエラー
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Internal server error",
		})
	}

	// 接続時の処理
	log.Println("New WebSocket connection established on port", h.service.Port)
	
	//接続によるメッセージ処理をゴルーチンに投げる
	go h.HandleMessages(conn)

	return nil
}