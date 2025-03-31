package repository_impl

import (
	"errors"
	"sync"

	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/domain/repository"
)

type WebsocketRepositoryImpl struct {
	clientData map[string]*entity.WebsocketClient
	mu         *sync.Mutex
}

func NewWebsocketRepositoryImpl() repository.IWebsocketRepository {
	return &WebsocketRepositoryImpl{
		clientData: make(map[string]*entity.WebsocketClient),
		mu:         &sync.Mutex{},
	}
}

func (wr *WebsocketRepositoryImpl) CreateClient(id string) error {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	if _, exists := wr.clientData[id]; exists {
		return errors.New("client already exists")
	}

	wr.clientData[id] = &entity.WebsocketClient{
		ID: id,
		SDP: "",
		Candidate: nil,
	}

	return nil
}

func (wr *WebsocketRepositoryImpl) SaveSDP(id string, sdp string) error {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	client, exists := wr.clientData[id]
	if !exists {
		return errors.New("client not found")
	}
	client.SDP = sdp

	return nil
}

func (wr *WebsocketRepositoryImpl) GetSDPByID(id string) (string, error) {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	client, exists := wr.clientData[id]
	if !exists || client.SDP == "" {
		return "", errors.New("SDP not found")
	}
	return client.SDP, nil
}

func (wr *WebsocketRepositoryImpl) SaveCandidate(id string, candidate string) error {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	client, exists := wr.clientData[id]
	if !exists {
		return errors.New("client not found")
	}
	client.Candidate = []string{candidate}
	
	return nil
}

func (wr *WebsocketRepositoryImpl) AddCandidate(id string, candidate string) error {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	client, exists := wr.clientData[id]
	if !exists || client.Candidate == nil {
		return errors.New("candidates not found")
	}

	client.Candidate = append(client.Candidate, candidate)
	return nil
}

func (wr *WebsocketRepositoryImpl) GetCandidatesByID(id string) ([]string, error) {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	client, exists := wr.clientData[id]
	if !exists {
		return nil, errors.New("candidates not found")
	}
	return client.Candidate, nil
}

func (wr *WebsocketRepositoryImpl) ExistsCandidateByID(id string) bool {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	client, exists := wr.clientData[id]
	return client.Candidate != nil && exists
}
