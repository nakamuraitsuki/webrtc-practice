package adapter_impl

import "github.com/gorilla/websocket"

type WebsocketConnectionAdapterImpl struct {
	conn *websocket.Conn
}

func NewWebsocketConnectionAdapterImpl(conn *websocket.Conn) *WebsocketConnectionAdapterImpl {
	return &WebsocketConnectionAdapterImpl{
		conn: conn,
	}
}

func (w *WebsocketConnectionAdapterImpl) ReadMessageFunc() (int, []byte, error) {
	return w.conn.ReadMessage()
}

func (w *WebsocketConnectionAdapterImpl) WriteMessageFunc(messageType int, data []byte) error {
	return w.conn.WriteMessage(messageType, data)
}

func (w *WebsocketConnectionAdapterImpl) CloseFunc() error {
	return w.conn.Close()
}