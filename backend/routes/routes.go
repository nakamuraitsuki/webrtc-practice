package routes

import (
	"example.com/webrtc-practice/config"
	"example.com/webrtc-practice/signaling"
	
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, cfg *config.Config, wsHandler *signaling.WebSocketHandler) {

	//webSocket関係のルーティング
	wsGroup := e.Group("/ws")
	wsHandler.Register(wsGroup)

}