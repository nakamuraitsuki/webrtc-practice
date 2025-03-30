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

func (w *WebSocketConnectionImpl) WriteMessage(messageType int, data []byte) error {
	return w.conn.WriteMessage(messageType, data)
}

func (w *WebSocketConnectionImpl) Close() error {
	return w.conn.Close()
}
