package entity

import "errors"

type Room struct {
	id      string // ""を許さない
	name    string // ""を許さない
	clients map[string]*WebsocketClient
}

func NewRoom(id string, name string) (*Room, error) {
	if id == "" {
		return nil, errors.New("id must not be empty")
	}
	if name == "" {
		return nil, errors.New("name must not be empty")
	}
	
	return &Room{
		id:      id,
		name:    name,
		clients: make(map[string]*WebsocketClient),
	}, nil
}
