package server

import (
	"github.com/labstack/echo/v4"
	"example.com/webrtc-practice/config"
)

func ServerStart() {
	cfg := config.LoadConfig()
	e := echo.New()

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}