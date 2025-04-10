// interface/factory/websocket_connection_factory.go
package factory

import (
	"net/http"

	"example.com/webrtc-practice/internal/domain/service"
	"example.com/webrtc-practice/internal/interface/adopter"
)

type WebsocketConnectionAdopterFactory interface {
	NewAdopter(conn service.WebSocketConnection) adopter.ConnAdopter
}

type WebsocketConnectionFactory interface {
	NewConnection(w http.ResponseWriter, r *http.Request) (service.WebSocketConnection, error)
}
