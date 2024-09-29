package usecases

import (
	"errors"
	"time"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/domain/repository"
	"github.com/dgrijalva/jwt-go"
)

type AuthUseCase struct {
	SecretKey        string
	MemberRepository repository.MemberRepository
}

func NewAuthUseCase(secretKey string, memberRepo repository.MemberRepository) *AuthUseCase {
	return &AuthUseCase{
		SecretKey:        secretKey,
		MemberRepository: memberRepo,
	}
}

// member 객체를 받도록 변경
func (uc *AuthUseCase) GenerateToken(member *domain.Member) (string, error) {
	claims := jwt.MapClaims{
		"user_id":       member.UserID,
		"is_admin":      member.IsAdmin,
		"member_number": member.MemberNumber,
		"exp":           time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.SecretKey))
}

func (uc *AuthUseCase) Authenticate(userID, password string) (*domain.Member, error) {
	member, err := uc.MemberRepository.GetByUserID(userID)
	if err != nil || member == nil {
		return nil, errors.New("Invalid credentials")
	}

	if !member.CheckPassword(password) {
		return nil, errors.New("Invalid credentials")
	}

	return member, nil
}
