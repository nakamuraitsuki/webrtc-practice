package usecase

import (
	"example.com/webrtc-practice/internal/repository"
)

type UserUsecase struct {
	repo repository.UserRepository
}
