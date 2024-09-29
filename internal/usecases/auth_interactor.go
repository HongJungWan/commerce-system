package usecases

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type AuthUseCase struct {
	SecretKey string
}

func NewAuthUseCase(secretKey string) *AuthUseCase {
	return &AuthUseCase{SecretKey: secretKey}
}

func (uc *AuthUseCase) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.SecretKey))
}
