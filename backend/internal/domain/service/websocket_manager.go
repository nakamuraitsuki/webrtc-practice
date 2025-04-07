package service

type WebSocketConnection interface {
	ReadMessage() (int, []byte, error)
	WriteMessage([]byte) error
	Close() error
}

type WebsocketManager interface {
	RegisterConnection(conn WebSocketConnection) error
	RegisterID(conn WebSocketConnection, id string)
	DeleteConnection(conn WebSocketConnection) error
	GetConnectionByID(id string) (WebSocketConnection, error)
	ExistsByID(id string) bool
}