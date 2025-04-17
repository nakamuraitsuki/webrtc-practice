package repository_impl

import (
	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/domain/repository"
)

type RoomRepositoryImpl struct {
	rooms map[string]*entity.Room
}

func NewRoomRepository() repository.IRoomRepository{
	return &RoomRepositoryImpl{
		rooms: make(map[string]*entity.Room),
	}
}

func (r *RoomRepositoryImpl) CreateRoom(name string) (string, error) {
	return "", nil
}

func (r *RoomRepositoryImpl) GetRoomByID(id string) (*entity.Room, error) {
	room, exists := r.rooms[id]
	if !exists {
		return nil, nil
	}
	return room, nil
}