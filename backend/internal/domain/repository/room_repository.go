package repository

import "example.com/webrtc-practice/internal/domain/entity"

type IRoomRepository interface {
	CreateRoom(name string) (string, error)      // 部屋を立て、部屋IDを返す
	GetRoomByID(id string) (*entity.Room, error) // 部屋IDから部屋を取得する
}
