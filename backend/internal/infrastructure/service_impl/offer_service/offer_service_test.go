package offerservice_test

import (
	"testing"

	offerservice "example.com/webrtc-practice/internal/infrastructure/service_impl/offer_service"
	"github.com/stretchr/testify/assert"
)

func TestOfferServiceImpl(t *testing.T) {
	service := offerservice.NewOfferService()

	t.Run("Initial state should be empty", func(t *testing.T) {
		assert.Equal(t, "", service.GetOffer())
		assert.False(t, service.IsOffer())
		assert.False(t, service.IsOfferID("someID"))
	})

	t.Run("SetOffer should store ID", func(t *testing.T) {
		service.SetOffer("user123")
		assert.Equal(t, "user123", service.GetOffer())
		assert.True(t, service.IsOffer())
		assert.True(t, service.IsOfferID("user123"))
		assert.False(t, service.IsOfferID("otherID"))
	})

	t.Run("ClearOffer should reset ID", func(t *testing.T) {
		service.ClearOffer()
		assert.Equal(t, "", service.GetOffer())
		assert.False(t, service.IsOffer())
		assert.False(t, service.IsOfferID("user123"))
	})
}
