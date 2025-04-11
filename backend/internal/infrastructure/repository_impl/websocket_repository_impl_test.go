package repository_impl_test

import (
	"testing"

	"example.com/webrtc-practice/internal/domain/repository"
	"example.com/webrtc-practice/internal/infrastructure/repository_impl"
	"github.com/stretchr/testify/assert"
)

func setupRepo() repository.IWebsocketRepository {
	return repository_impl.NewWebsocketRepository()
}

func TestCreateAndDeleteClient(t *testing.T) {
	repo := setupRepo()
	id := "testID"

	t.Run("CreateClient", func(t *testing.T) {
		err := repo.CreateClient(id)
		assert.NoError(t, err)
	})

	t.Run("CreateClient_Duplicate", func(t *testing.T) {
		err := repo.CreateClient(id)
		assert.Error(t, err)
	})

	t.Run("DeleteClient", func(t *testing.T) {
		err := repo.DeleteClient(id)
		assert.NoError(t, err)
	})

	t.Run("DeleteClient_NotFound", func(t *testing.T) {
		err := repo.DeleteClient(id)
		assert.Error(t, err)
	})
}

func TestSDP(t *testing.T) {
	repo := setupRepo()
	id := "user1"
	_ = repo.CreateClient(id)

	t.Run("SaveSDP", func(t *testing.T) {
		err := repo.SaveSDP(id, "testSDP")
		assert.NoError(t, err)
	})

	t.Run("GetSDP", func(t *testing.T) {
		sdp, err := repo.GetSDPByID(id)
		assert.NoError(t, err)
		assert.Equal(t, "testSDP", sdp)
	})

	t.Run("GetSDPNotFound", func(t *testing.T) {
		_, err := repo.GetSDPByID("nonexistent")
		assert.Error(t, err)
	})
}

func TestCandidate(t *testing.T) {
	repo := setupRepo()
	id := "user1"
	_ = repo.CreateClient(id)

	t.Run("SaveCandidate", func(t *testing.T) {
		err := repo.SaveCandidate(id, []string{"cand1"})
		assert.NoError(t, err)
	})

	t.Run("AddCandidate", func(t *testing.T) {
		err := repo.AddCandidate(id, []string{"cand2"})
		assert.NoError(t, err)
	})

	t.Run("GetCandidates", func(t *testing.T) {
		candidates, err := repo.GetCandidatesByID(id)
		assert.NoError(t, err)
		assert.Equal(t, []string{"cand1", "cand2"}, candidates)
	})

	t.Run("ExistsCandidate", func(t *testing.T) {
		exists := repo.ExistsCandidateByID(id)
		assert.True(t, exists)
	})

	t.Run("SaveCandidate_InvalidID", func(t *testing.T) {
		err := repo.SaveCandidate("invalid", []string{"cand"})
		assert.Error(t, err)
	})

	t.Run("AddCandidate_InvalidID", func(t *testing.T) {
		err := repo.AddCandidate("invalid", []string{"cand"})
		assert.Error(t, err)
	})

	t.Run("GetCandidates_InvalidID", func(t *testing.T) {
		_, err := repo.GetCandidatesByID("invalid")
		assert.Error(t, err)
	})

	t.Run("ExistsCandidate_InvalidID", func(t *testing.T) {
		exists := repo.ExistsCandidateByID("invalid")
		assert.False(t, exists)
	})
}
