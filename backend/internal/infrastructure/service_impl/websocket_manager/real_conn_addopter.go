package websocketmanager

import "github.com/gorilla/websocket"

type RealConnAdopter interface {
	ReadMessageFunc() (int, []byte, error)
	WriteMessageFunc(int, []byte) error
	CloseFunc() error
}

type RealConnAdopterImpl struct {
	conn *websocket.Conn
}

func NewRealConnAdopter(conn *websocket.Conn) RealConnAdopter {
	return &RealConnAdopterImpl{
		conn: conn,
	}
}

func (r *RealConnAdopterImpl) ReadMessageFunc() (int, []byte, error) {
	return r.conn.ReadMessage()
}

func (r *RealConnAdopterImpl) WriteMessageFunc(messageType int, data []byte) error {
	return r.conn.WriteMessage(messageType, data)
}

func (r *RealConnAdopterImpl) CloseFunc() error {
	return r.conn.Close()
}