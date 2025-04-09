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
		err := manager.RegisterID(mockConn, "testID")
		assert.NoError(t, err)
	})
	t.Run("未登録コネクションに対するID登録", func(t *testing.T) {
		mockConn2 := mock_service.NewMockWebSocketConnection(ctrl)
		err := manager.RegisterID(mockConn2, "testID")
		assert.Error(t, err)
		assert.EqualError(t, err, "client not registered")
	})
}

func TestDeleteConnection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConn1 := mock_service.NewMockWebSocketConnection(ctrl)
	testID := "testID"
	mockConn2 := mock_service.NewMockWebSocketConnection(ctrl)
	mockConn3 := mock_service.NewMockWebSocketConnection(ctrl)
	manager := websocketmanager.NewWebsocketManager()

	err := manager.RegisterConnection(mockConn1)
	if err != nil {
		t.Fatalf("failed to register connection: %v", err)
	}
	err = manager.RegisterID(mockConn1, testID)
	if err != nil {
		t.Fatalf("failed to register ID: %v", err)
	}
	err = manager.RegisterConnection(mockConn2)
	if err != nil {
		t.Fatalf("failed to register connection: %v", err)
	}

	t.Run("ID登録済みコネクションの削除", func(t *testing.T) {
		err := manager.DeleteConnection(mockConn1)
		assert.NoError(t, err)
	})

	t.Run("ID未登録コネクションの削除", func(t *testing.T) {
		err := manager.DeleteConnection(mockConn2)
		assert.NoError(t, err)
	})

	t.Run("登録されていないコネクションはエラーを返す", func(t *testing.T) {
		err := manager.DeleteConnection(mockConn3)
		assert.Error(t, err)
		assert.EqualError(t, err, "client not found")
	})
}

func TestGetConnectionByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConn := mock_service.NewMockWebSocketConnection(ctrl)
	testID := "testID"
	manager := websocketmanager.NewWebsocketManager()
	
	err := manager.RegisterConnection(mockConn)
	if err != nil {
		t.Fatalf("failed to register connection: %v", err)
	}
	err = manager.RegisterID(mockConn, testID)
	if err != nil {
		t.Fatalf("failed to register ID: %v", err)
	}

	t.Run("ID登録済みコネクションの取得", func(t *testing.T) {
		conn, err := manager.GetConnectionByID(testID)
		assert.NoError(t, err)
		assert.Equal(t, mockConn, conn)
	})

	t.Run("未登録IDに対するエラー", func(t *testing.T) {
		_, err := manager.GetConnectionByID("nonexistentID")
		assert.Error(t, err)
		assert.EqualError(t, err, "client not found")
	})
}

