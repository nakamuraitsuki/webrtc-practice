package repository

import (
	"example.com/webrtc-practice/internal/domain/entity"
)

type CreateUserParams struct {
	Name     string
	Email    string
	PasswdHash string
}

type UserRepository interface {
	CreateUser(params CreateUserParams) error
	GetAllUsers() ([]*entity.User, error)
	GetUserByID(id int) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	UpdateUser(user *entity.User) error
	DeleteUser(id int) error
}
