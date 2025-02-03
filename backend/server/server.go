package server

import (
	"example.com/webrtc-practice/config"
	"example.com/webrtc-practice/routes"
	"github.com/labstack/echo/v4"
)

func ServerStart() {
	cfg := config.LoadConfig()
	e := echo.New()

	routes.SetupRoutes(e, cfg)

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}