package usecase

import (
	"example.com/webrtc-practice/internal/domain/service"
	"example.com/webrtc-practice/internal/domain/repository"
)

type UserUsecase struct {
	repo repository.UserRepository
	hasyer service.Hasher
}
