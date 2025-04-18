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

func NewWebsocketRepository() repository.IWebsocketRepository {
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

	wr.clientData[id] = entity.NewWebsocketClient(id,"",nil)

	return nil
}

func (wr *WebsocketRepositoryImpl) DeleteClient(id string) error {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()
	if _, exists := wr.clientData[id]; !exists {
		return errors.New("client not found")
	}
	delete(wr.clientData, id)
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
	client.SetSDP(sdp)

	return nil
}

func (wr *WebsocketRepositoryImpl) GetSDPByID(id string) (string, error) {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	client, exists := wr.clientData[id]
	if !exists || client.GetSDP() == "" {
		return "", errors.New("SDP not found")
	}
	return client.GetSDP(), nil
}

func (wr *WebsocketRepositoryImpl) SaveCandidate(id string, candidate []string) error {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	client, exists := wr.clientData[id]
	if !exists {
		return errors.New("client not found")
	}
	client.SetCandidate(append(client.GetCandidate(), candidate...))
	
	return nil
}

func (wr *WebsocketRepositoryImpl) AddCandidate(id string, candidate []string) error {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	client, exists := wr.clientData[id]
	if !exists {
		return errors.New("client not found")
	}
	client.SetCandidate(append(client.GetCandidate(), candidate...))
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
	return client.GetCandidate(), nil
}

func (wr *WebsocketRepositoryImpl) ExistsCandidateByID(id string) bool {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	client, exists := wr.clientData[id]
	// existsを先に調べ、nullポインタへのアクセスを防ぐ
	return exists && client.GetCandidate() != nil
}
