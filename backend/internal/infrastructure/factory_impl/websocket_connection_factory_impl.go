package factory_impl

import (
	"net/http"

	"example.com/webrtc-practice/internal/domain/service"
	"example.com/webrtc-practice/internal/infrastructure/adapter_impl"
	websocketmanager "example.com/webrtc-practice/internal/infrastructure/service_impl/websocket_manager"
	websocketupgrader "example.com/webrtc-practice/internal/infrastructure/service_impl/websocket_upgrader"
)

type WebsocketConnectionFactoryImpl struct {
	Upgrader websocketupgrader.WebsocketUpgraderInterface
}



func NewWebsocketConnectionFactoryImpl(upgrader websocketupgrader.WebsocketUpgraderInterface) *WebsocketConnectionFactoryImpl {
	return &WebsocketConnectionFactoryImpl{
		Upgrader: upgrader,
	}
}

func (w *WebsocketConnectionFactoryImpl) NewConnection(wr http.ResponseWriter, r *http.Request) (service.WebSocketConnection, error) {
	conn, err := w.Upgrader.Upgrade(wr, r, nil)
	if err != nil {
		return nil, err
	}

	client := adapter_impl.NewWebsocketConnectionAdapterImpl(conn)

	return websocketmanager.NewWebsocketConnection(client), nil
}
