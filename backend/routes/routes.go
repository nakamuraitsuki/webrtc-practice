package routes

import (
	"net/http"

	"example.com/webrtc-practice/config"
	"example.com/webrtc-practice/signaling"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, cfg *config.Config) {

	//webSocket関係のルーティング
	wsGroup := e.Group("/ws")

	// WebSocketHandlerをインスタンス化
	webSocketService := &signaling.WebSocketService{Port: cfg.Port}
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	webSocketHandler := signaling.NewWebSocketHandler(webSocketService, upgrader)
	webSocketHandler.Register(wsGroup)
}