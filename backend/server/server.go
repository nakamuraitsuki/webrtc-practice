package server

import (
	"example.com/webrtc-practice/config"
	"example.com/webrtc-practice/internal/handler"
	"example.com/webrtc-practice/internal/infrastructure/repository_impl"
	"example.com/webrtc-practice/internal/infrastructure/repository_impl/sqlite3"
	"example.com/webrtc-practice/internal/infrastructure/service_impl/hasher"
	"example.com/webrtc-practice/internal/infrastructure/service_impl/jwt"
	offerservice "example.com/webrtc-practice/internal/infrastructure/service_impl/offer_service"
	websocketbroadcast "example.com/webrtc-practice/internal/infrastructure/service_impl/websocket_broadcast"
	websocketmanager "example.com/webrtc-practice/internal/infrastructure/service_impl/websocket_manager"
	"example.com/webrtc-practice/routes"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ServerStart(cfg *config.Config, db *sqlx.DB) {
	e := echo.New()

	// ミドルウェアの設定
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},                                              // すべてのオリジンを許可
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},       // 許可するHTTPメソッド
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization}, // 許可するHTTPヘッダー
	}))

	// ユーザーハンドラの初期化
	userRepository := sqlite3.NewUserRepository(db)
	hasher := hasher.NewBcryptHasher()
	tokenService := jwt.NewJWTService(cfg.SecretKey, cfg.TokenExpiry)
	userHandler := handler.NewUserHandler(userRepository, hasher, tokenService)

	// WebSocketハンドラの初期化
	websocketRepository := repository_impl.NewWebsocketRepository()
	websocketManager := websocketmanager.NewWebsocketManager()
	websocketBroadcast := websocketbroadcast.NewBroadcast()
	websocketOfferService := offerservice.NewOfferService()

	websocketHandler := handler.NewWebsocketHandler(
		websocketRepository,
		websocketManager,
		websocketBroadcast,
		websocketOfferService,
	)

	routes.SetupRoutes(e, cfg, userHandler, websocketHandler)

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
