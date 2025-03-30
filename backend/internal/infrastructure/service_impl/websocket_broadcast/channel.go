package websocketbroadcast

import "example.com/webrtc-practice/internal/domain/service"

// TODO ： Message型をやり取りするようにする
type Broadcast struct {
	broadcast chan []byte
}

func NewBroadcast() service.WebSocketBroadcastService {
	return &Broadcast{
		broadcast: make(chan []byte),
	}
}

func (b *Broadcast) Send(message []byte) {
	b.broadcast <- message
}

func (b *Broadcast) Receive() []byte {
	return <-b.broadcast
}
