package handler_test

import (
	"testing"
	"time"

	"example.com/webrtc-practice/internal/handler"
	mock_service_impl "example.com/webrtc-practice/mocks/infrastructure/service_impl"
	mock_usecase "example.com/webrtc-practice/mocks/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewWebsocketHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// モック作成
	mockUsecase := mock_usecase.NewMockIWebsocketUsecaseInterface(ctrl)
	mockUpgrader := mock_service_impl.NewMockWebsocketUpgraderInterface(ctrl)

	wsHandler := handler.NewWebsocketHandler(mockUsecase, mockUpgrader)

	// ProcessMessageが1回は呼ばれることを期待（HandleMessagesの中で呼ばれる）
	mockUsecase.EXPECT().ProcessMessage().Times(1)

	assert.NotNil(t, wsHandler)

	// goroutineの実行を少し待つ
	time.Sleep(10 * time.Millisecond)
}
