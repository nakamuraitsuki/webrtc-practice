package server

import (
	"net/http"

	"example.com/webrtc-practice/config"
	"example.com/webrtc-practice/routes"
	"example.com/webrtc-practice/signaling"
	
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

func ServerStart() {
	cfg := config.LoadConfig()
	e := echo.New()

	//WebSocketの依存関係
	websocketService := &signaling.WebSocketService{Port: cfg.Port}
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	websocketHandler := signaling.NewWebSocketHandler(websocketService, upgrader)

	routes.SetupRoutes(e, cfg, websocketHandler)

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}