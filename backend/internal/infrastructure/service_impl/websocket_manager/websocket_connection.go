package websocketmanager

import (
	"encoding/json"

	"example.com/webrtc-practice/internal/domain/entity"
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

func (w *WebSocketConnectionImpl) ReadMessage() (int, entity.Message, error) {
	messageType, messagebyte, err := w.conn.ReadMessage()
	if err != nil {
		return 0, entity.Message{}, err
	}

	var message entity.Message
	err = json.Unmarshal(messagebyte, &message)
	if err != nil {
		return 0, entity.Message{}, err
	}
	

	return messageType, message, nil
}

func (w *WebSocketConnectionImpl) WriteMessage(data entity.Message) error {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return w.conn.WriteMessage(websocket.TextMessage, dataByte)
}

func (w *WebSocketConnectionImpl) Close() error {
	return w.conn.Close()
}
