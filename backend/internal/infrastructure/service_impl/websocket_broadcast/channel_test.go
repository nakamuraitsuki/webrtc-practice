package websocketbroadcast_test

import (
	"testing"
	"time"

	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/domain/service"
	websocketbroadcast "example.com/webrtc-practice/internal/infrastructure/service_impl/websocket_broadcast"
	"github.com/stretchr/testify/assert"
)

func TestBroadcast(t *testing.T) {
	b := websocketbroadcast.NewBroadcast()

	t.Run("Send and Receive Message", func(t *testing.T) {
		expected, _ := entity.NewMessage("123", "connection", "sdp", []string{"candidate", "candidate2"}, "456")
		// チャネルがブロックしないように別ゴルーチンで受信
		go func() {
			b.Send(*expected)
		}()

		select {
		case <-time.After(1 * time.Second):
			t.Fatal("Receive timed out, message was not received")
		case received := <-getChan(b):
			assert.Equal(t, expected, received)
		}
	})
}

// helper関数: interfaceでラップされた Broadcast からチャネルを取り出せないため、
// Receive() を呼び出して受信するゴルーチンを作る
func getChan(b service.WebSocketBroadcastService) <-chan entity.Message {
	ch := make(chan entity.Message)
	go func() {
		ch <- b.Receive()
	}()
	return ch
}
