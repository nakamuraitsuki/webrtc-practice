package entity

type Message struct {
	ID        string  `json:"id"` // ユーザーIDを期待する
	Type      string  `json:"type"`
	SDP       string `json:"sdp"`
	Candidate []string `json:"candidate"`
	TargetID  string `json:"target_id"`
}

type WebsocketClient struct {
	ID        string   `json:"id"`
	SDP       string  `json:"sdp"`
	Candidate []string  `json:"candidate"`
}

func NewMessage(id string, messageType string, sdp string, candidate []string, targetID string) *Message {
	return &Message{
		ID:        id,
		Type:      messageType,
		SDP:       sdp,
		Candidate: candidate,
		TargetID:  targetID,
	}
}

func NewWebsocketClient(id, sdp string, candidate []string) *WebsocketClient {
	return &WebsocketClient{
		ID:        id,
		SDP:       sdp,
		Candidate: candidate,
	}
}
