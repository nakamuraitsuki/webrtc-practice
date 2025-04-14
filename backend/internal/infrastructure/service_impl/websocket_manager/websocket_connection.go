package websocketmanager

import (
	"encoding/json"

	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/domain/service"
	"example.com/webrtc-practice/internal/infrastructure/dto"
	"example.com/webrtc-practice/internal/interface/adapter"
	"github.com/gorilla/websocket"
)

type WebSocketConnectionImpl struct {
	conn adapter.ConnAdapter
}

func NewWebsocketConnection(conn adapter.ConnAdapter) service.WebSocketConnection {
	return &WebSocketConnectionImpl{
		conn: conn,
	}
}

func (w *WebSocketConnectionImpl) ReadMessage() (int, entity.Message, error) {
	messageType, messagebyte, err := w.conn.ReadMessageFunc()
	if err != nil {
		return 0, entity.Message{}, err
	}

	var message dto.WebsocketMessageDTO
	err = json.Unmarshal(messagebyte, &message)
	if err != nil {
		return 0, entity.Message{}, err
	}
	
	messageEntity := message.ToEntity()
	return messageType, *messageEntity, nil
}

func (w *WebSocketConnectionImpl) WriteMessage(data entity.Message) error {
	dataDTO := dto.WebsocketMessageDTO{}
	dataDTO.FromEntity(&data)
	dataByte, err := json.Marshal(dataDTO)
	if err != nil {
		return err
	}
	return w.conn.WriteMessageFunc(websocket.TextMessage, dataByte)
}

func (w *WebSocketConnectionImpl) Close() error {
	return w.conn.CloseFunc()
}
