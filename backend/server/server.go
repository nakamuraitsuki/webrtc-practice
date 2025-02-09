package server

import (
	"example.com/webrtc-practice/config"
	"example.com/webrtc-practice/internal/handler"
	"example.com/webrtc-practice/internal/infrastructure/repository/sqlite3"
	"example.com/webrtc-practice/internal/infrastructure/service/hasher"
	"example.com/webrtc-practice/internal/infrastructure/service/jwt"
	"example.com/webrtc-practice/routes"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func ServerStart(cfg *config.Config, db *sqlx.DB) {
	e := echo.New()

	// ユーザーハンドラーの初期化
	userRepository := sqlite3.NewUserRepository(db)
	hasher := hasher.NewBcryptHasher()
	tokenService := jwt.NewJWTService(cfg.SecretKey, cfg.TokenExpiry)
	userHandler := handler.NewUserHandler(userRepository, hasher, tokenService)

	routes.SetupRoutes(e, cfg, userHandler)

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
