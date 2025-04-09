package websocketmanager_test

import (
	"testing"

	websocketmanager "example.com/webrtc-practice/internal/infrastructure/service_impl/websocket_manager"
	mock_service "example.com/webrtc-practice/mocks/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewWebsocketManager(t *testing.T) {
	mnager := websocketmanager.NewWebsocketManager()

	assert.NotNil(t, mnager)
}

func TestRegisterConnection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConn := mock_service.NewMockWebSocketConnection(ctrl)
	manager := websocketmanager.NewWebsocketManager()

	t.Run("初回コネクション登録", func(t *testing.T) {
		err := manager.RegisterConnection(mockConn)
		assert.NoError(t, err)
	})

	t.Run("重複コネクション登録", func(t *testing.T) {
		err := manager.RegisterConnection(mockConn)
		assert.Error(t, err)
		assert.EqualError(t, err, "client already registered")
	})
}

func TestRegisterID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConn := mock_service.NewMockWebSocketConnection(ctrl)
	manager := websocketmanager.NewWebsocketManager()

	err := manager.RegisterConnection(mockConn)
	if err != nil {
		t.Fatalf("failed to register connection: %v", err)
	}

	t.Run("既存コネクションに対するID登録", func(t *testing.T) {
		manager.RegisterID(mockConn, "testID")
	}
}

