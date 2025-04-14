package entity_test

import (
	"testing"

	"example.com/webrtc-practice/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestMessage_Getters(t *testing.T) {
	message := entity.NewMessage("id", "type", "sdp", []string{"c1"}, "target")

	assert.Equal(t, "id", message.GetID())
	assert.Equal(t, "type", message.GetType())
	assert.Equal(t, "sdp", message.GetSDP())
	assert.Equal(t, []string{"c1"}, message.GetCandidate())
	assert.Equal(t, "target", message.GetTargetID())
}

func TestMessage_Setters(t *testing.T) {
	message := entity.NewMessage("id", "type", "sdp", []string{"c1"}, "target")

	message.SetID("new_id")
	message.SetType("new_type")
	message.SetSDP("new_sdp")
	message.SetCandidate([]string{"new_c1"})
	message.SetTargetID("new_target")

	assert.Equal(t, "new_id", message.GetID())
	assert.Equal(t, "new_type", message.GetType())
	assert.Equal(t, "new_sdp", message.GetSDP())
	assert.Equal(t, []string{"new_c1"}, message.GetCandidate())
	assert.Equal(t, "new_target", message.GetTargetID())
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

	client.SetID("new_user")
	client.SetSDP("new_sdp")
	client.SetCandidate([]string{"new_candidate"})

	assert.Equal(t, "new_user", client.GetID())
	assert.Equal(t, "new_sdp", client.GetSDP())
	assert.Equal(t, []string{"new_candidate"}, client.GetCandidate())
}