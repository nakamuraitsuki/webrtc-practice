package jwt

import (
	"time"
	"github.com/dgrijalva/jwt-go"

	"example.com/webrtc-practice/internal/domain/service"
)

type JWTSercice struct {
	SecretKey   string
	TokenExpiry time.Duration
}

func NewJWTService(secretKey string, tokenExpiry time.Duration) service.TokenService {
	return &JWTSercice{
		SecretKey:   secretKey,
		TokenExpiry: tokenExpiry,
	}
}

func (js *JWTSercice) GenerateToken(userID int) (string, error) {
	// JWTペイロード設定
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(js.TokenExpiry).Unix(),
	}

	// JWT作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// JWT署名
	signedToken, err := token.SignedString([]byte(js.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (js *JWTSercice) ValidateToken(token string) (int, error) {
	// JWT検証
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(js.SecretKey), nil
	})
	if err != nil {
		return 0, err
	}

	// ユーザID取得
	userID := int(claims["user_id"].(float64))

	return userID, nil
}
