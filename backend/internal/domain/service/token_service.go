package service

type TokenService interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(token string) (int, error)
}