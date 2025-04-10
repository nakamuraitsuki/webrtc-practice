package websocketmanager

import (
	"encoding/json"

	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/domain/service"
	"example.com/webrtc-practice/internal/interface/adopter"
	"github.com/gorilla/websocket"
)

type WebSocketConnectionImpl struct {
	conn adopter.ConnAdopter
}

func NewWebsocketConnection(conn adopter.ConnAdopter) service.WebSocketConnection {
	return &WebSocketConnectionImpl{
		conn: conn,
	}
}

func (w *WebSocketConnectionImpl) ReadMessage() (int, entity.Message, error) {
	messageType, messagebyte, err := w.conn.ReadMessageFunc()
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
	return w.conn.WriteMessageFunc(websocket.TextMessage, dataByte)
}

func (w *WebSocketConnectionImpl) Close() error {
	return w.conn.CloseFunc()
}
