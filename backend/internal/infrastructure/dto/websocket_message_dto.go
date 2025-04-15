package dto

import "example.com/webrtc-practice/internal/domain/entity"

type WebsocketMessageDTO struct {
	ID        string   `json:"id"`
	Type      string   `json:"type"`
	SDP       string   `json:"sdp"`
	Candidate []string `json:"candidate"`
	TargetID  string   `json:"target_id"`
}

func (w *WebsocketMessageDTO) ToEntity() *entity.Message {
	res, err := entity.NewMessage(
		w.ID,
		w.Type,
		w.SDP,
		w.Candidate,
		w.TargetID,
	)
	if err != nil {
		return nil
	}
	return res
}

func (w *WebsocketMessageDTO) FromEntity(msg *entity.Message) {
	w.ID = msg.GetID()
	w.Type = msg.GetType()
	w.SDP = msg.GetSDP()
	w.Candidate = msg.GetCandidate()
	w.TargetID = msg.GetTargetID()
}