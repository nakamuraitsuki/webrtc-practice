package adopter_impl

import "github.com/gorilla/websocket"

type WebsocketConnectionAdopterImpl struct {
	conn *websocket.Conn
}

func NewWebsocketConnectionAdopterImpl(conn *websocket.Conn) *WebsocketConnectionAdopterImpl {
	return &WebsocketConnectionAdopterImpl{
		conn: conn,
	}
}

func (w *WebsocketConnectionAdopterImpl) ReadMessageFunc() (int, []byte, error) {
	return w.conn.ReadMessage()
}

func (w *WebsocketConnectionAdopterImpl) WriteMessageFunc(messageType int, data []byte) error {
	return w.conn.WriteMessage(messageType, data)
}

func (w *WebsocketConnectionAdopterImpl) CloseFunc() error {
	return w.conn.Close()
}