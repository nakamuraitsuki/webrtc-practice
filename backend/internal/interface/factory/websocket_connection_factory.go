// interface/factory/websocket_connection_factory.go
package factory

import (
	"net/http"

	"example.com/webrtc-practice/internal/domain/service"
	"example.com/webrtc-practice/internal/interface/adapter"
)

type WebsocketConnectionAdapterFactory interface {
	NewAdapter(conn service.WebSocketConnection) adapter.ConnAdapter
}

type WebsocketConnectionFactory interface {
	NewConnection(w http.ResponseWriter, r *http.Request) (service.WebSocketConnection, error)
}
