package entity

type Message struct {
	id        string // ユーザーIDを期待する
	msgType   string
	sdp       string
	candidate []string
	targetID  string
}

type WebsocketClient struct {
	id        string   
	sdp       string   
	candidate []string 
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

func (m *Message) SetCandidate(candidate []string) {
	m.candidate = candidate
}

func (m *Message) SetTargetID(targetID string) {
	m.targetID = targetID
}

func NewWebsocketClient(id, sdp string, candidate []string) *WebsocketClient {
	return &WebsocketClient{
		id:        id,
		sdp:       sdp,
		candidate: candidate,
	}
}

func (w *WebsocketClient) GetID() string {
	return w.id
}

func (w *WebsocketClient) GetSDP() string {
	return w.sdp
}

func (w *WebsocketClient) GetCandidate() []string {
	return w.candidate
}

func (w *WebsocketClient) SetSDP(sdp string) {
	w.sdp = sdp
}

func (w *WebsocketClient) SetCandidate(candidate []string) {
	w.candidate = candidate
}
