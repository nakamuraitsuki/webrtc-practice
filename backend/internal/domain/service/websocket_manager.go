package service

import "example.com/webrtc-practice/internal/domain/entity"

type WebSocketConnection interface {
	ReadMessage() (int, entity.Message, error)
	WriteMessage(entity.Message) error
	Close() error
}

type WebsocketManager interface {
	RegisterConnection(conn WebSocketConnection) error
	RegisterID(conn WebSocketConnection, id string)
	DeleteConnection(conn WebSocketConnection) error
	GetConnectionByID(id string) (WebSocketConnection, error)
	ExistsByID(id string) bool
}