package entity

import "fmt"

const (
	// msgTypeの定義
	MESSAGE_TYPE_CONNECT   = "connect"
	MESSAGE_TYPE_OFFLINE   = "offline"
	MESSAGE_TYPE_SDP       = "sdp"
	MESSAGE_TYPE_CANDIDATE = "candidate"
)

type MessageType string

type Message struct {
	id        string // ユーザーIDを期待する
	msgType   MessageType
	sdp       string
	candidate []string
	targetID  string
}

type WebsocketClient struct {
	id        string
	sdp       string
	candidate []string
}

func NewMessageType(messageType string) (MessageType, error) {
	mt := MessageType(messageType)
	switch mt {
	case MESSAGE_TYPE_CONNECT, MESSAGE_TYPE_OFFLINE, MESSAGE_TYPE_SDP, MESSAGE_TYPE_CANDIDATE:
		return mt, nil
	default:
		return "", fmt.Errorf("invalid message type: %s", messageType)
	}
}

func NewMessage(id string, messageType string, sdp string, candidate []string, targetID string) (*Message, error) {
	mt, err := NewMessageType(messageType)
	if err != nil {
		return nil, err
	}

	return &Message{
		id:        id,
		msgType:   mt,
		sdp:       sdp,
		candidate: candidate,
		targetID:  targetID,
	}, nil
}

func (m Message) GetID() string {
	return m.id
}

func (m Message) GetType() string {
	return string(m.msgType)
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
