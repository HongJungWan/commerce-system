package repository_test

import (
	"testing"
	"time"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/infrastructure/repository"
	"github.com/HongJungWan/commerce-system/test/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestMemberRepositoryImpl_Create_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	repo := repository.NewMemberRepository(db)
	member := &domain.Member{
		ID:           12345,
		MemberNumber: "M12345",
		AccountId:    "testuser",
		Password:     "password123",
		NickName:     "Test User",
		Email:        "testuser@example.com",
	}

	// When
	err := repo.Create(member)

	// Then
	assert.NoError(t, err)
	assert.NotZero(t, member.ID)
}

func TestMemberRepositoryImpl_Create_Failure_DuplicateUserID(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	repo := repository.NewMemberRepository(db)
	member1 := &domain.Member{
		ID:           12345,
		MemberNumber: "M12345",
		AccountId:    "duplicateuser",
		Password:     "password123",
		NickName:     "User One",
		Email:        "user1@example.com",
	}
	member2 := &domain.Member{
		ID:           12346,
		MemberNumber: "M12346",
		AccountId:    "duplicateuser",
		Password:     "password123",
		NickName:     "User Two",
		Email:        "user2@example.com",
	}
	_ = repo.Create(member1)

	// When
	err := repo.Create(member2)

	// Then
	assert.Error(t, err)
}

func TestMemberRepositoryImpl_GetByID_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	repo := repository.NewMemberRepository(db)
	member := &domain.Member{
		ID:           12345,
		MemberNumber: "M12345",
		AccountId:    "testuser",
		Password:     "password123",
		NickName:     "Test User",
		Email:        "testuser@example.com",
	}
	_ = repo.Create(member)

	// When
	retrievedMember, err := repo.GetByID(member.ID)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, member.AccountId, retrievedMember.AccountId)
}

func TestMemberRepositoryImpl_GetByID_Failure_NotFound(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	repo := repository.NewMemberRepository(db)

	// When
	retrievedMember, err := repo.GetByID(999)

	// Then
	assert.Error(t, err)
	assert.Nil(t, retrievedMember)
}

func TestMemberRepositoryImpl_GetByUserID_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	repo := repository.NewMemberRepository(db)
	member := &domain.Member{
		ID:           12345,
		MemberNumber: "M12345",
		AccountId:    "testuser",
		Password:     "password123",
		NickName:     "Test User",
		Email:        "testuser@example.com",
	}
	_ = repo.Create(member)

	// When
	retrievedMember, err := repo.GetByUserName("testuser")

	// Then
	assert.NoError(t, err)
	assert.Equal(t, member.MemberNumber, retrievedMember.MemberNumber)
}

func TestMemberRepositoryImpl_GetByUserID_Failure_NotFound(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	repo := repository.NewMemberRepository(db)

	// When
	retrievedMember, err := repo.GetByUserName("nonexistent")

	// Then
	assert.Error(t, err)
	assert.Nil(t, retrievedMember)
}

func TestMemberRepositoryImpl_Update_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	repo := repository.NewMemberRepository(db)
	member := &domain.Member{
		ID:           12345,
		MemberNumber: "M12345",
		AccountId:    "testuser",
		Password:     "password123",
		NickName:     "Old Name",
		Email:        "old@example.com",
	}
	_ = repo.Create(member)

	// When
	member.NickName = "New Name"
	err := repo.Update(member)

	// Then
	assert.NoError(t, err)

	// Verify
	updatedMember, _ := repo.GetByID(member.ID)
	assert.Equal(t, "New Name", updatedMember.NickName)
}

func TestMemberRepositoryImpl_Delete_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	repo := repository.NewMemberRepository(db)
	member := &domain.Member{
		ID:           12345,
		MemberNumber: "M12345",
		AccountId:    "testuser",
		Password:     "password123",
		NickName:     "Test User",
		Email:        "testuser@example.com",
	}
	_ = repo.Create(member)

	// When
	err := repo.Delete(member.ID)

	// Then
	assert.NoError(t, err)

	// Verify
	deletedMember, err := repo.GetByID(member.ID)
	assert.NoError(t, err)
	assert.True(t, deletedMember.IsWithdrawn)
}

func TestMemberRepositoryImpl_GetAll_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	repo := repository.NewMemberRepository(db)
	member1 := &domain.Member{
		ID:           12345,
		MemberNumber: "M12345",
		AccountId:    "user1",
		Password:     "password123",
		NickName:     "User One",
		Email:        "user1@example.com",
	}
	member2 := &domain.Member{
		ID:           12346,
		MemberNumber: "M12346",
		AccountId:    "user2",
		Password:     "password123",
		NickName:     "User Two",
		Email:        "user2@example.com",
	}
	_ = repo.Create(member1)
	_ = repo.Create(member2)

	// When
	members, err := repo.GetAll()

	// Then
	assert.NoError(t, err)
	assert.Len(t, members, 2)
}

func TestMemberRepositoryImpl_GetStatsByMonth_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	repo := repository.NewMemberRepository(db)

	// 특정 월에 가입한 회원 생성
	member1 := &domain.Member{
		ID:           12345,
		MemberNumber: "M12345",
		AccountId:    "user1",
		Password:     "password123",
		NickName:     "User One",
		Email:        "user1@example.com",
		CreatedAt:    time.Date(2024, 9, 10, 0, 0, 0, 0, time.UTC),
	}
	member2 := &domain.Member{
		ID:           12346,
		MemberNumber: "M12346",
		AccountId:    "user2",
		Password:     "password123",
		NickName:     "User Two",
		Email:        "user2@example.com",
		CreatedAt:    time.Date(2024, 9, 15, 0, 0, 0, 0, time.UTC),
	}
	member3 := &domain.Member{
		ID:           12347,
		MemberNumber: "M12347",
		AccountId:    "user3",
		Password:     "password123",
		NickName:     "User Three",
		Email:        "user3@example.com",
		CreatedAt:    time.Date(2024, 8, 20, 0, 0, 0, 0, time.UTC),
	}

	_ = repo.Create(member1)
	_ = repo.Create(member2)
	_ = repo.Create(member3)

	// 회원 탈퇴
	_ = repo.Delete(member2.ID)

	// 삭제된 회원의 WithdrawnAt을 테스트 월에 포함되도록 설정
	member2Fetched, _ := repo.GetByID(member2.ID)
	withdrawnAt := time.Date(2024, 9, 20, 0, 0, 0, 0, time.UTC)
	member2Fetched.WithdrawnAt = &withdrawnAt
	_ = repo.Update(member2Fetched)

	// When
	joined, deleted, err := repo.GetStatsByMonth("2024-09")

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 2, joined)  // 수정된 기대 값
	assert.Equal(t, 1, deleted) // 수정된 기대 값
}

func TestMemberRepositoryImpl_GetStatsByMonth_Failure_InvalidMonth(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	repo := repository.NewMemberRepository(db)

	// When
	joined, deleted, err := repo.GetStatsByMonth("invalid-month")

	// Then
	assert.Error(t, err)
	assert.Equal(t, 0, joined)
	assert.Equal(t, 0, deleted)
}
