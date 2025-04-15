package entity_test

import (
	"testing"

	"example.com/webrtc-practice/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestMessage_Getters(t *testing.T) {
	message, _ := entity.NewMessage("id", "connect", "sdp", []string{"c1"}, "target")

	assert.Equal(t, "id", message.GetID())
	assert.Equal(t, "connect", message.GetType())
	assert.Equal(t, "sdp", message.GetSDP())
	assert.Equal(t, []string{"c1"}, message.GetCandidate())
	assert.Equal(t, "target", message.GetTargetID())
}

func TestWebsocketClient_Getters(t *testing.T) {
	id := "user123"
	sdp := "sdp data"
	candidate := []string{"candidate1", "candidate2"}

	client := entity.NewWebsocketClient(id, sdp, candidate)

	assert.Equal(t, id, client.GetID())
	assert.Equal(t, sdp, client.GetSDP())
	assert.Equal(t, candidate, client.GetCandidate())
}

func TestWebsocketClient_Setters(t *testing.T) {
	id := "user123"
	sdp := "sdp data"
	candidate := []string{"candidate1", "candidate2"}

	client := entity.NewWebsocketClient(id, sdp, candidate)

	client.SetSDP("new_sdp")
	client.SetCandidate([]string{"new_candidate"})

	assert.Equal(t, "new_sdp", client.GetSDP())
	assert.Equal(t, []string{"new_candidate"}, client.GetCandidate())
}