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
		"account_id":    member.AccountId,
		"is_admin":      member.IsAdmin == true,
		"member_number": member.MemberNumber,
		"exp":           time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.SecretKey))
}

func (uc *AuthUseCase) Authenticate(userName, password string) (*domain.Member, error) {
	member, err := uc.MemberRepository.GetByAccountId(userName)
	if err != nil || member == nil {
		return nil, errors.New("Invalid credentials")
	}

	if !member.CheckPassword(password) {
		return nil, errors.New("Invalid credentials")
	}

	return member, nil
}
