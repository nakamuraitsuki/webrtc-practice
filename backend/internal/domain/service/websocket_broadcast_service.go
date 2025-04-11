package service

import "example.com/webrtc-practice/internal/domain/entity"

type WebSocketBroadcastService interface {
	Send(message entity.Message)
	Receive() entity.Message
}