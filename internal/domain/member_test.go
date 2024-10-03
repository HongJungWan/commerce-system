package domain_test

import (
	"testing"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestMember_SetPassword_Success(t *testing.T) {
	// Given
	member := &domain.Member{}

	// When
	err := member.AssignPassword("password123")

	// Then
	assert.NoError(t, err)
	assert.NotEmpty(t, member.Password)
}

func TestMember_CheckPassword_CorrectPassword(t *testing.T) {
	// Given
	member := &domain.Member{}
	_ = member.AssignPassword("password123")

	// When
	isValid := member.CheckPassword("password123")

	// Then
	assert.True(t, isValid)
}

func TestMember_CheckPassword_IncorrectPassword(t *testing.T) {
	// Given
	member := &domain.Member{}
	_ = member.AssignPassword("password123")

	// When
	isValid := member.CheckPassword("wrongpassword")

	// Then
	assert.False(t, isValid)
}

func TestMember_Validate_Success(t *testing.T) {
	// Given
	member := &domain.Member{
		ID:           12345,
		MemberNumber: "M12345",
		AccountId:    "testuser",
		Password:     "password123",
		NickName:     "Test User",
		Email:        "testuser@example.com",
	}

	// When
	err := member.Validate()

	// Then
	assert.NoError(t, err)
}

func TestMember_Validate_Failure_MissingFields(t *testing.T) {
	// Given
	member := &domain.Member{}

	// When
	err := member.Validate()

	// Then
	assert.Error(t, err)
	assert.Equal(t, "회원번호가 누락되었습니다.", err.Error())
}
