package repository_impl

import (
	"errors"
	"sync"
)

type WebsocketRepositoryImpl struct {
	sdpData       map[string]string
	candidateData map[string][]string
	mu            *sync.Mutex
}

func NewWebsocketRepositoryImpl() *WebsocketRepositoryImpl {
	return &WebsocketRepositoryImpl{
		sdpData:       make(map[string]string),
		candidateData: make(map[string][]string),
		mu:            &sync.Mutex{},
	}
}

func (wr *WebsocketRepositoryImpl) SaveSDP(id string, sdp string) error {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	wr.sdpData[id] = sdp
	return nil
}

func (wr *WebsocketRepositoryImpl) GetSDPByID(id string) (string, error) {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	sdp, exists := wr.sdpData[id]
	if !exists {
		return "", errors.New("SDP not found")
	}
	return sdp, nil
}

func (wr *WebsocketRepositoryImpl) SaveCandidate(id string, candidate string) error {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	wr.candidateData[id] = []string{candidate}
	return nil
}

func (wr *WebsocketRepositoryImpl) AddCandidate(id string, candidate string) error {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	if _, exists := wr.candidateData[id]; !exists {
		return errors.New("candidates not found")
	}

	wr.candidateData[id] = append(wr.candidateData[id], candidate)
	return nil
}

func (wr *WebsocketRepositoryImpl) GetCandidatesByID(id string) ([]string, error) {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	candidates, exists := wr.candidateData[id]
	if !exists {
		return nil, errors.New("candidates not found")
	}
	return candidates, nil
}

func (wr *WebsocketRepositoryImpl) ExistsCandidateByID(id string) bool {
	// ミューテーションロックを使用して、同時アクセスを防止
	wr.mu.Lock()
	defer wr.mu.Unlock()

	_, exists := wr.candidateData[id]
	return exists
}
