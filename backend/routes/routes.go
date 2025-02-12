package routes

import (
	"example.com/webrtc-practice/config"
	"example.com/webrtc-practice/internal/handler"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, cfg *config.Config, userHandler handler.UserHandler) {

	userGroup := e.Group("/api/user")
	userHandler.Register(userGroup)
}
