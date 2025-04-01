package entity_test

import (
	"testing"

	"example.com/webrtc-practice/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewMessage(t *testing.T) {
	id := "user123"
	messageType := "offer"
	sdp := "sdp data"
	candidate := []string{"candidate1", "candidate2"}
	targetID := "target456"

	message := entity.NewMessage(id, messageType, sdp, candidate, targetID)

	assert.Equal(t, id, message.ID)
	assert.Equal(t, messageType, message.Type)
	assert.Equal(t, sdp, message.SDP)
	assert.Equal(t, candidate, message.Candidate)
	assert.Equal(t, targetID, message.TargetID)
}

func TestNewWebsocketClient(t *testing.T) {
	id := "user123"
	sdp := "sdp data"
	candidate := []string{"candidate1", "candidate2"}

	client := entity.NewWebsocketClient(id, sdp, candidate)

	assert.Equal(t, id, client.ID)
	assert.Equal(t, sdp, client.SDP)
	assert.Equal(t, candidate, client.Candidate)
}
