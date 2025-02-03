package routes

import (
	"example.com/webrtc-practice/config"
	"example.com/webrtc-practice/signaling"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, cfg *config.Config) {

	//webSocket関係のルーティング
	wsGroup := e.Group("/ws")

	// WebSocketHandlerをインスタンス化
	webSocketService := &signaling.WebSocketService{Port: cfg.Port}
	webSocketHandler := signaling.NewWebSocketHandler(webSocketService)
	webSocketHandler.Register(wsGroup)
}