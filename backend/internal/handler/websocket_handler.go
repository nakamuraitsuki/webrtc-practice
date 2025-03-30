package handler

import (
	"net/http"

	"example.com/webrtc-practice/internal/usecase"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{}
)

type WebsocketHandler struct {
	Usecase usecase.IWebsocketUsecase
}

func NewWebsocketHandler(
) WebsocketHandler {
	h := WebsocketHandler{
		Usecase: usecase.NewWebsocketUsecase(),
	}

	// WebSocketメッセージ処理のゴルーチンを起動
	go h.HandleMessages()

	return h
}

func (h *WebsocketHandler) Register(g *echo.Group) {
	g.GET("/", h.HandleWebSocket)
}

// WebSocket接続
func (h *WebsocketHandler) HandleWebSocket(c echo.Context) error {
	// リクエストをコネクションにアップグレード
	conn, _ := upgrader.Upgrade(c.Response().Writer, c.Request(), nil) // error ignored for sake of simplicity
	defer conn.Close()

	err := h.Usecase.RegisterClient(conn)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "client already registered"})
	}

	go h.Usecase.ListenForMessages(conn)

	return c.JSON(http.StatusOK, map[string]string{"message": "success"})
}

// メッセージ処理の呼び出し
func (h *WebsocketHandler) HandleMessages() {
	h.Usecase.ProcessMessage()
}
