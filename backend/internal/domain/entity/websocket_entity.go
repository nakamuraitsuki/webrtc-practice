package entity

type Message struct {
	ID        string  `json:"id"` // ユーザーIDを期待する
	Type      string  `json:"type"`
	SDP       string `json:"sdp"`
	Candidate string `json:"candidate"`
	TargetID  string `json:"target_id"`
}

type WebsocketClient struct {
	ID        string   `json:"id"`
	SDP       string  `json:"sdp"`
	Candidate []string  `json:"candidate"`
}
