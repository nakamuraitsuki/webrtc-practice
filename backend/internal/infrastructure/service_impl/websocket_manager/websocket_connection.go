package websocketmanager

import (
	"example.com/webrtc-practice/internal/domain/service"
	"github.com/gorilla/websocket"
)

type WebSocketConnectionImpl struct {
	conn *websocket.Conn
}

func NewWebsocketConnection(conn *websocket.Conn) service.WebSocketConnection {
	return &WebSocketConnectionImpl{
		conn: conn,
	}
}

func (w *WebSocketConnectionImpl) ReadMessage() (int, []byte, error) {
	return w.conn.ReadMessage()
}

func (w *WebSocketConnectionImpl) WriteMessage(data []byte) error {
	return w.conn.WriteMessage(websocket.TextMessage, data)
}

func (w *WebSocketConnectionImpl) Close() error {
	return w.conn.Close()
}
