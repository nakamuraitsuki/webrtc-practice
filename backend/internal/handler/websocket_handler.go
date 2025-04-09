package handler

import (
	"net/http"

	websocketmanager "example.com/webrtc-practice/internal/infrastructure/service_impl/websocket_manager"
	websocketupgrader "example.com/webrtc-practice/internal/infrastructure/service_impl/websocket_upgrader"
	"example.com/webrtc-practice/internal/usecase"
	"github.com/labstack/echo/v4"
)

type WebsocketHandler struct {
	Usecase  usecase.IWebsocketUsecaseInterface
	Upgrader websocketupgrader.WebsocketUpgraderInterface
}

func NewWebsocketHandler(
	usecase usecase.IWebsocketUsecaseInterface,
	upgrader websocketupgrader.WebsocketUpgraderInterface,
) WebsocketHandler {
	h := WebsocketHandler{
		Usecase:  usecase,
		Upgrader: upgrader,
	}

	// WebSocketメッセージ処理のゴルーチンを起動
	go h.HandleMessages()

	return h
}

func (h *WebsocketHandler) Register(g *echo.Group) {
	g.GET("/", h.HandleWebsocket)
}

// WebSocket接続
func (h *WebsocketHandler) HandleWebsocket(c echo.Context) error {
	// リクエストをコネクションにアップグレード
	conn, _ := h.Upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	defer conn.Close()

	connAdopter := websocketmanager.NewRealConnAdopter(conn)
	client := websocketmanager.NewWebsocketConnection(connAdopter)

	err := h.Usecase.RegisterClient(client)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "client already registered"})
	}

	go h.Usecase.ListenForMessages(client)

	return c.JSON(http.StatusOK, map[string]string{"message": "success"})
}

// メッセージ処理の呼び出し
func (h *WebsocketHandler) HandleMessages() {
	h.Usecase.ProcessMessage()
}
