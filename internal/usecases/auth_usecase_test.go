package usecases_test

import (
	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/infrastructure/repository"
	"github.com/HongJungWan/commerce-system/internal/usecases"
	"github.com/HongJungWan/commerce-system/test/fixtures"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthUseCase_Authenticate_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)

	member := &domain.Member{
		UserID:       "testuser",
		Password:     "password123",
		Name:         "Test User",
		Email:        "testuser@example.com",
		MemberNumber: "M12345",
	}
	_ = member.SetPassword(member.Password)
	_ = memberRepo.Create(member)

	// When
	authenticatedMember, err := authUseCase.Authenticate("testuser", "password123")

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "testuser", authenticatedMember.UserID)
}

func TestAuthUseCase_Authenticate_Failure_WrongPassword(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)

	member := &domain.Member{
		UserID:       "testuser",
		Password:     "password123",
		Name:         "Test User",
		Email:        "testuser@example.com",
		MemberNumber: "M12345",
	}
	_ = member.SetPassword(member.Password)
	_ = memberRepo.Create(member)

	// When
	authenticatedMember, err := authUseCase.Authenticate("testuser", "wrongpassword")

	// Then
	assert.Error(t, err)
	assert.Nil(t, authenticatedMember)
	assert.Equal(t, "Invalid credentials", err.Error())
}

func TestAuthUseCase_Authenticate_Failure_UserNotFound(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)

	// When
	authenticatedMember, err := authUseCase.Authenticate("nonexistent", "password123")

	// Then
	assert.Error(t, err)
	assert.Nil(t, authenticatedMember)
	assert.Equal(t, "Invalid credentials", err.Error())
}

func TestAuthUseCase_GenerateToken_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)

	member := &domain.Member{
		UserID:       "testuser",
		IsAdmin:      true,
		MemberNumber: "M12345",
	}

	// When
	token, err := authUseCase.GenerateToken(member)

	// Then
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
