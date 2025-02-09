package usecase

import (
	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/domain/repository"
	"example.com/webrtc-practice/internal/domain/service"
)

type UserUsecase struct {
	repo         repository.IUserRepository
	hasher       service.Hasher
	tokenService service.TokenService
}

func NewUserUsecase(repo repository.IUserRepository, hasher service.Hasher) *UserUsecase {
	return &UserUsecase{
		repo:   repo,
		hasher: hasher,
	}
}

func (u *UserUsecase) RegisterUser(name, email, password string) (*entity.User, error) {
	hashedPassword, err := u.hasher.HashPassword(password)
	if err != nil {
		return nil, err
	}

	res, err := u.repo.CreateUser(repository.CreateUserParams{
		Name:       name,
		Email:      email,
		PasswdHash: hashedPassword,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *UserUsecase) AuthenticateUser(email, password string) (string, error) {
	user, err := u.repo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	ok, err := u.hasher.ComparePassword(user.GetPasswdHash(), password)
	if err != nil {
		return "", err
	}

	if !ok {
		return "", err
	}

	return u.tokenService.GenerateToken(user.GetID())
}
