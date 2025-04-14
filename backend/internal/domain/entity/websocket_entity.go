package entity

type Message struct {
	id        string    // ユーザーIDを期待する
	msgType   string   
	sdp       string   
	candidate []string 
	targetID  string   
}

type WebsocketClient struct {
	ID        string   `json:"id"`
	SDP       string   `json:"sdp"`
	Candidate []string `json:"candidate"`
}

func NewMessage(id string, messageType string, sdp string, candidate []string, targetID string) *Message {
	return &Message{
		id:        id,
		msgType:   messageType,
		sdp:       sdp,
		candidate: candidate,
		targetID:  targetID,
	}
}

func (m Message) GetID() string {
	return m.id
}

func (m Message) GetType() string {
	return m.msgType
}

func (m Message) GetSDP() string {
	return m.sdp
}

func (m Message) GetCandidate() []string {
	return m.candidate
}

func (m Message) GetTargetID() string {
	return m.targetID
}

func (m *Message) SetID(id string) {
	m.id = id
}

func (m *Message) SetType(messageType string) {
	m.msgType = messageType
}

func (m *Message) SetSDP(sdp string) {
	m.sdp = sdp
}

func (m *Message) SetCandidate(candidate []string) {
	m.candidate = candidate
}

func (m *Message) SetTargetID(targetID string) {
	m.targetID = targetID
}

func NewWebsocketClient(id, sdp string, candidate []string) *WebsocketClient {
	return &WebsocketClient{
		ID:        id,
		SDP:       sdp,
		Candidate: candidate,
	}
}
