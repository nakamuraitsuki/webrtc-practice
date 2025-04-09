package websocketupgrader

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type WebsocketUpgraderInterface interface {
	// Upgrade upgrades the HTTP connection to a WebSocket connection.
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error)
}

type WebsocketUpgrader struct {
	Upgrader websocket.Upgrader
}

func NewWebsocketUpgrader() WebsocketUpgraderInterface {
	return &WebsocketUpgrader{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			Subprotocols: []string{"json"},
		},
	}
}

func (w *WebsocketUpgrader) Upgrade(wr http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error) {
	conn, err := w.Upgrader.Upgrade(wr, r, responseHeader)
	if err != nil {
		return nil, err
	}
	return conn, nil
}