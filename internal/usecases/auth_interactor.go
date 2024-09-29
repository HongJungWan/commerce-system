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

func (uc *AuthUseCase) GenerateToken(member *domain.Member) (string, error) {
	claims := jwt.MapClaims{
		"user_id":       member.UserID,
		"is_admin":      member.IsAdmin == true, // FIXME: 명시적으로 bool 타입으로 설정, 추후 수정 필요
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
