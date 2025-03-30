package repository

type IWebsocketRepository interface {
	CreateClient(id string) error
	SaveSDP(id string, sdp string) error
	GetSDPByID(id string) (string, error)
	SaveCandidate(id string, candidate string) error
	AddCandidate(id string, candidate string) error
	ExistsCandidateByID(id string) bool
	GetCandidatesByID(id string) ([]string, error)
	DeleteSDP(id string) error
}