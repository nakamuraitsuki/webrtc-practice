package service

type WebSocketBroadcastService interface {
	Send(message []byte)
	Receive() []byte
}