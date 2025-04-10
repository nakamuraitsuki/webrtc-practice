package handler_test

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"example.com/webrtc-practice/internal/handler"
	mock_factory "example.com/webrtc-practice/mocks/interface/factory"
	mock_service "example.com/webrtc-practice/mocks/service"
	mock_usecase "example.com/webrtc-practice/mocks/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewWebsocketHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// モック作成
	mockUsecase := mock_usecase.NewMockIWebsocketUsecaseInterface(ctrl)
	mockFactory := mock_factory.NewMockWebsocketConnectionFactory(ctrl)

	wsHandler := handler.NewWebsocketHandler(mockUsecase, mockFactory)

	// ProcessMessageが1回は呼ばれることを期待（HandleMessagesの中で呼ばれる）
	mockUsecase.EXPECT().ProcessMessage().Times(1)

	assert.NotNil(t, wsHandler)

	// goroutineの実行を少し待つ
	time.Sleep(10 * time.Millisecond)
}

func TestHandleWebsocket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// モック作成
	mockUsecase := mock_usecase.NewMockIWebsocketUsecaseInterface(ctrl)
	mockFactory := mock_factory.NewMockWebsocketConnectionFactory(ctrl)

	dummyContext, _ := createTestContext(echo.GET, "/ws", "")

	wsHandler := handler.NewWebsocketHandler(mockUsecase, mockFactory)

	// NewWebsocketHandler呼び出し時のゴルーチンの実行を待つ
	mockUsecase.EXPECT().ProcessMessage().Times(1)
	time.Sleep(10 * time.Millisecond)

	// WebSocket接続のモックを作成
	mockConn := mock_service.NewMockWebSocketConnection(ctrl)

	// WebSocket接続のアップグレードをモック
	mockFactory.EXPECT().NewConnection(gomock.Any(), gomock.Any()).Return(mockConn, nil).Times(1)

	// RegisterClientが1回呼ばれることを期待
	mockUsecase.EXPECT().RegisterClient(mockConn).Return(nil).Times(1)

	// ListenForMessagesが1回呼ばれることを期待
	mockUsecase.EXPECT().ListenForMessages(mockConn).Times(1)

	mockConn.EXPECT().Close().Times(1)

	// HandleWebsocketメソッドをテスト
	err := wsHandler.HandleWebsocket(dummyContext)
	assert.NoError(t, err)

	// goroutineの実行を少し待つ
	time.Sleep(10 * time.Millisecond)
}

func createTestContext(method, path, body string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()

	// リクエスト作成
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// レスポンス記録用
	rec := httptest.NewRecorder()

	// Context生成
	c := e.NewContext(req, rec)

	return c, rec
}
