package usecase

import (
	"example.com/webrtc-practice/internal/domain/repository"
	"example.com/webrtc-practice/internal/domain/service"
)

type UserUsecase struct {
	repo   repository.UserRepository
	hasher service.Hasher
}

func NewUserUsecase(repo repository.UserRepository, hasher service.Hasher) *UserUsecase {
	return &UserUsecase{
		repo:   repo,
		hasher: hasher,
	}
}

func (u *UserUsecase) RegisterUser(name, email, password string) error {
	hashedPassword, err := u.hasher.HashPassword(password)
	if err != nil {
		return err
	}

	return u.repo.CreateUser(repository.CreateUserParams{
		Name:      name,
		Email:     email,
		PasswdHash: hashedPassword,
	})
}