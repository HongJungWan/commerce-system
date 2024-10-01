package usecases_test

import (
	"testing"
	"time"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/infrastructure/repository"
	"github.com/HongJungWan/commerce-system/internal/usecases"
	"github.com/HongJungWan/commerce-system/test/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestMemberInteractor_Register_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	interactor := usecases.NewMemberInteractor(memberRepo)

	member := &domain.Member{
		Username:       "newuser",
		HashedPassword: "password123",
		FullName:       "New User",
		Email:          "newuser@example.com",
		MemberNumber:   "M12345",
	}

	// When
	err := interactor.Register(member)

	// Then
	assert.NoError(t, err)
	retrievedMember, _ := memberRepo.GetByUserName("newuser")
	assert.NotNil(t, retrievedMember)
}

func TestMemberInteractor_Register_Failure_DuplicateUserID(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	interactor := usecases.NewMemberInteractor(memberRepo)

	existingMember := &domain.Member{
		Username:       "duplicateuser",
		HashedPassword: "password123",
		FullName:       "Existing User",
		Email:          "existing@example.com",
		MemberNumber:   "M12345",
	}
	_ = memberRepo.Create(existingMember)

	newMember := &domain.Member{
		Username:       "duplicateuser",
		HashedPassword: "password456",
		FullName:       "New User",
		Email:          "new@example.com",
		MemberNumber:   "M12346",
	}

	// When
	err := interactor.Register(newMember)

	// Then
	assert.Error(t, err)
	assert.Equal(t, "이미 존재하는 사용자 ID입니다.", err.Error())
}

func TestMemberInteractor_GetByUserID_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	interactor := usecases.NewMemberInteractor(memberRepo)

	member := &domain.Member{
		Username:       "testuser",
		HashedPassword: "password123",
		FullName:       "Test User",
		Email:          "testuser@example.com",
		MemberNumber:   "M12345",
	}
	_ = memberRepo.Create(member)

	// When
	retrievedMember, err := interactor.GetByUserName("testuser")

	// Then
	assert.NoError(t, err)
	assert.Equal(t, member.FullName, retrievedMember.FullName)
}

func TestMemberInteractor_GetByUserID_Failure_NotFound(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	interactor := usecases.NewMemberInteractor(memberRepo)

	// When
	retrievedMember, err := interactor.GetByUserName("nonexistent")

	// Then
	assert.Error(t, err)
	assert.Nil(t, retrievedMember)
}

func TestMemberInteractor_UpdateByUserID_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	interactor := usecases.NewMemberInteractor(memberRepo)

	member := &domain.Member{
		Username:       "testuser",
		HashedPassword: "password123",
		FullName:       "Old Name",
		Email:          "old@example.com",
		MemberNumber:   "M12345",
	}
	_ = memberRepo.Create(member)

	updateData := &domain.Member{
		FullName: "New Name",
		Email:    "new@example.com",
	}

	// When
	err := interactor.UpdateByUserName("testuser", updateData)

	// Then
	assert.NoError(t, err)
	updatedMember, _ := memberRepo.GetByUserName("testuser")
	assert.Equal(t, "New Name", updatedMember.FullName)
	assert.Equal(t, "new@example.com", updatedMember.Email)
}

func TestMemberInteractor_UpdateByUserID_Failure_UserNotFound(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	interactor := usecases.NewMemberInteractor(memberRepo)

	updateData := &domain.Member{
		FullName: "New Name",
		Email:    "new@example.com",
	}

	// When
	err := interactor.UpdateByUserName("nonexistent", updateData)

	// Then
	assert.Error(t, err)
}

func TestMemberInteractor_DeleteByUserID_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	interactor := usecases.NewMemberInteractor(memberRepo)

	member := &domain.Member{
		Username:       "testuser",
		HashedPassword: "password123",
		FullName:       "Test User",
		Email:          "testuser@example.com",
		MemberNumber:   "M12345",
	}
	_ = memberRepo.Create(member)

	// When
	err := interactor.DeleteByUserName("testuser")

	// Then
	assert.NoError(t, err)
	deletedMember, err := memberRepo.GetByUserName("testuser")
	assert.Error(t, err)
	assert.Nil(t, deletedMember)
}

func TestMemberInteractor_DeleteByUserID_Failure_UserNotFound(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	interactor := usecases.NewMemberInteractor(memberRepo)

	// When
	err := interactor.DeleteByUserName("nonexistent")

	// Then
	assert.Error(t, err)
}

func TestMemberInteractor_GetAll_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	interactor := usecases.NewMemberInteractor(memberRepo)

	member1 := &domain.Member{
		Username:       "user1",
		HashedPassword: "password123",
		FullName:       "User One",
		Email:          "user1@example.com",
		MemberNumber:   "M12345",
	}
	member2 := &domain.Member{
		Username:       "user2",
		HashedPassword: "password123",
		FullName:       "User Two",
		Email:          "user2@example.com",
		MemberNumber:   "M12346",
	}
	_ = memberRepo.Create(member1)
	_ = memberRepo.Create(member2)

	// When
	members, err := interactor.GetAll()

	// Then
	assert.NoError(t, err)
	assert.Len(t, members, 2)
}

func TestMemberInteractor_GetStatsByMonth_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	interactor := usecases.NewMemberInteractor(memberRepo)

	member1 := &domain.Member{
		Username:       "user1",
		HashedPassword: "password123",
		FullName:       "User One",
		Email:          "user1@example.com",
		MemberNumber:   "M12345",
		CreatedAt:      time.Date(2024, 9, 10, 0, 0, 0, 0, time.UTC),
	}
	member2 := &domain.Member{
		Username:       "user2",
		HashedPassword: "password123",
		FullName:       "User Two",
		Email:          "user2@example.com",
		MemberNumber:   "M12346",
		CreatedAt:      time.Date(2024, 9, 15, 0, 0, 0, 0, time.UTC),
	}
	_ = memberRepo.Create(member1)
	_ = memberRepo.Create(member2)
	_ = memberRepo.Delete(member2.ID)

	// When
	joined, deleted, err := interactor.GetStatsByMonth("2024-09")

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 1, joined)
	assert.Equal(t, 0, deleted)
}

func TestMemberInteractor_GetStatsByMonth_Failure_InvalidMonth(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	interactor := usecases.NewMemberInteractor(memberRepo)

	// When
	joined, deleted, err := interactor.GetStatsByMonth("invalid-month")

	// Then
	assert.Error(t, err)
	assert.Equal(t, 0, joined)
	assert.Equal(t, 0, deleted)
}
