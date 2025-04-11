package websocketbroadcast

import (
	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/domain/service"
)

type Broadcast struct {
	broadcast chan entity.Message
}

func NewBroadcast() service.WebSocketBroadcastService {
	return &Broadcast{
		broadcast: make(chan entity.Message),
	}
}

func (b *Broadcast) Send(message entity.Message) {
	b.broadcast <- message
}

func (b *Broadcast) Receive() entity.Message {
	return <-b.broadcast
}
