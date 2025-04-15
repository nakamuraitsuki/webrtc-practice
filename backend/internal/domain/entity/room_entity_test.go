package entity_test

import (
	"testing"

	"example.com/webrtc-practice/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewRoom(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		room, err := entity.NewRoom("room1", "test room")
		assert.NoError(t, err)
		assert.NotNil(t, room)
	})
	t.Run("異常系: IDが空文字", func(t *testing.T) {
		room, err := entity.NewRoom("", "test room")
		assert.Error(t, err)
		assert.Nil(t, room)
	})
	t.Run("異常系: 名前が空文字", func(t *testing.T) {
		room, err := entity.NewRoom("room1", "")
		assert.Error(t, err)
		assert.Nil(t, room)
	})
}
