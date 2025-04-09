package websocketmanager_test

import (
	"testing"


)

func TestWebsocketManager_RegisterConnection_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// モック作成
	mockConn := mock_service.NewMockWebSocketConnection(ctrl)

	// WebSocketManagerインスタンス作成
	manager := websocket_manager.NewWebsocketManager()

	// RegisterConnectionのテスト
	err := manager.RegisterConnection(mockConn)
	require.NoError(t, err)

	// 同じ接続で再登録を試みるとエラーになるかテスト
	err = manager.RegisterConnection(mockConn)
	assert.Error(t, err, "client already registered")
}

func TestWebsocketManager_RegisterID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// モック作成
	mockConn := mock.NewMockWebSocketConnection(ctrl)
	manager := websocket_manager.NewWebsocketManager()

	// RegisterConnectionを呼び出して登録
	err := manager.RegisterConnection(mockConn)
	require.NoError(t, err)

	// RegisterIDでIDを登録
	manager.RegisterID(mockConn, "client123")

	// IDを取得できるかテスト
	conn, err := manager.GetConnectionByID("client123")
	require.NoError(t, err)
	assert.Equal(t, mockConn, conn)
}

func TestWebsocketManager_DeleteConnection_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// モック作成
	mockConn := mock.NewMockWebSocketConnection(ctrl)
	manager := websocket_manager.NewWebsocketManager()

	// RegisterConnectionを呼び出して登録
	err := manager.RegisterConnection(mockConn)
	require.NoError(t, err)

	// RegisterIDでIDを登録
	manager.RegisterID(mockConn, "client123")

	// DeleteConnectionで削除
	err = manager.DeleteConnection(mockConn)
	require.NoError(t, err)

	// 削除後、IDでの接続取得を試みる
	_, err = manager.GetConnectionByID("client123")
	assert.Error(t, err, "client not found")
}

func TestWebsocketManager_ExistsByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// モック作成
	mockConn := mock.NewMockWebSocketConnection(ctrl)
	manager := websocket_manager.NewWebsocketManager()

	// RegisterConnectionを呼び出して登録
	err := manager.RegisterConnection(mockConn)
	require.NoError(t, err)

	// RegisterIDでIDを登録
	manager.RegisterID(mockConn, "client123")

	// IDが存在するか確認
	exists := manager.ExistsByID("client123")
	assert.True(t, exists)

	// 存在しないIDで確認
	exists = manager.ExistsByID("nonexistent")
	assert.False(t, exists)
}

func TestWebsocketManager_DeleteConnection_ClientNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// モック作成
	mockConn := mock.NewMockWebSocketConnection(ctrl)
	manager := NewWebsocketManager()

	// 存在しない接続を削除しようとする
	err := manager.DeleteConnection(mockConn)
	assert.Error(t, err, "client not found")
}
