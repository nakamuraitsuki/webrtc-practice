package test

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"example.com/webrtc-practice/internal/infrastructure/factory_impl"
	"example.com/webrtc-practice/internal/infrastructure/repository_impl"
	offerservice "example.com/webrtc-practice/internal/infrastructure/service_impl/offer_service"
	websocketbroadcast "example.com/webrtc-practice/internal/infrastructure/service_impl/websocket_broadcast"
	websocketmanager "example.com/webrtc-practice/internal/infrastructure/service_impl/websocket_manager"
	websocketupgrader "example.com/webrtc-practice/internal/infrastructure/service_impl/websocket_upgrader"
	"example.com/webrtc-practice/internal/interface/handler"
	"example.com/webrtc-practice/internal/usecase"
)

func TestWebsocketHandler_E2E_Real(t *testing.T) {
	e := echo.New()

	// 本物の依存関係を組み立てる（Factory, Repository, Usecase）
	wsRepo := repository_impl.NewWebsocketRepository()
	wsManager := websocketmanager.NewWebsocketManager()
	wsBr := websocketbroadcast.NewBroadcast()
	wsO := offerservice.NewOfferService()
	wsUsecase := usecase.NewWebsocketUsecase(wsRepo, wsManager, wsBr, wsO)
	upgrader := websocketupgrader.NewWebsocketUpgrader()
	wsFactory := factory_impl.NewWebsocketConnectionFactoryImpl(upgrader)

	// ハンドラを作ってルーティング
	wsHandler := handler.NewWebsocketHandler(wsUsecase, wsFactory)
	wsHandler.Register(e.Group("/ws"))

	// サーバー起動（httptest経由）
	server := httptest.NewServer(e)
	defer server.Close()

	// WebSocket用URLへ変換
	wsURL := "ws" + server.URL[len("http"):] + "/ws/"

	// === テスト1: 接続 ===
	t.Run("クライアントがwebsocket接続できる", func(t *testing.T) {
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		assert.NoError(t, err)
		defer conn.Close()
	})

	// === テスト2: メッセージ送信 ===
	t.Run("クライアントが相互接続できる", func(t *testing.T) {
		conn1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		assert.NoError(t, err)
		defer conn1.Close()

		conn2, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		assert.NoError(t, err)
		defer conn2.Close()

		// 少し待機して処理が安定するようにする
		time.Sleep(100 * time.Millisecond)

		testConnectMessage1 := `{
			"id": "client1",
			"type": "connect",
			"sdp": "",
			"candidate": [],
			"target_id": ""
		}`

		testOfferMessage1 := `{
			"id": "client1",
			"type": "offer",
			"sdp": "sdp1",
			"candidate": [],
			"target_id": ""
		}`

		testConnectMessage2 := `{
			"id": "client2",
			"type": "connect",
			"sdp": "",
			"candidate": [],
			"target_id": ""
		}`

		err = conn1.WriteMessage(websocket.TextMessage, []byte(testConnectMessage1))
		assert.NoError(t, err)

		conn1.SetReadDeadline(time.Now().Add(2 * time.Second))
		// client1がフィードバックメッセージを受け取れているか確認
		_, msg1, err := conn1.ReadMessage()
		assert.NoError(t, err)
		t.Logf("client1 received: %s", msg1)
		assert.Contains(t, string(msg1), `"type":"offer"`)
		assert.Contains(t, string(msg1), `"id":"client1"`)

		err = conn1.WriteMessage(websocket.TextMessage, []byte(testOfferMessage1))
		assert.NoError(t, err)

		time.Sleep(100 * time.Millisecond)
		err = conn2.WriteMessage(websocket.TextMessage, []byte(testConnectMessage2))
		assert.NoError(t, err)

		conn2.SetReadDeadline(time.Now().Add(2 * time.Second))
		// client2がclient1からのofferを受け取れているか確認
		_, msg2, err := conn2.ReadMessage()
		assert.NoError(t, err)
		t.Logf("client2 received: %s", msg2)
		assert.Contains(t, string(msg2), `"type":"offer"`)
		assert.Contains(t, string(msg2), `"id":"client1"`)
		assert.Contains(t, string(msg2), `"target_id":"client2"`)
	})
	// TODO: 一連のやり取りをすべてテストするべき
}
