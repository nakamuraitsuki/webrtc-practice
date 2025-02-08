package server

import (
	"example.com/webrtc-practice/config"

	"github.com/labstack/echo/v4"
)

func ServerStart() {
	cfg := config.LoadConfig()
	e := echo.New()

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
