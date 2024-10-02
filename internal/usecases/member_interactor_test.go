package usecases_test

import (
	"testing"
	"time"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/infrastructure/repository"
	"github.com/HongJungWan/commerce-system/internal/interfaces/dto/request"
	"github.com/HongJungWan/commerce-system/internal/usecases"
	"github.com/HongJungWan/commerce-system/test/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestMemberInteractor_Register_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	interactor := usecases.NewMemberInteractor(memberRepo, authUseCase)

	req := &request.CreateMemberRequest{
		AccountId: "newuser",
		Password:  "password123",
		NickName:  "New User",
		Email:     "newuser@example.com",
	}

	// When
	responseData, err := interactor.Register(req)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "회원 가입이 완료되었습니다.", responseData.Message)
	retrievedMember, _ := memberRepo.GetByUserName("newuser")
	assert.NotNil(t, retrievedMember)
}

func TestMemberInteractor_Register_Failure_DuplicateUserID(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	interactor := usecases.NewMemberInteractor(memberRepo, authUseCase)

	existingMember := &domain.Member{
		MemberNumber: "M12345",
		Username:     "duplicateuser",
		FullName:     "Existing User",
		Email:        "existing@example.com",
	}
	existingMember.AssignPassword("password123")
	_ = memberRepo.Create(existingMember)

	req := &request.CreateMemberRequest{
		AccountId: "duplicateuser",
		Password:  "password456",
		NickName:  "New User",
		Email:     "new@example.com",
	}

	// When
	_, err := interactor.Register(req)

	// Then
	assert.Error(t, err)
	assert.Equal(t, "이미 존재하는 사용자 ID입니다.", err.Error())
}

func TestMemberInteractor_GetMyInfo_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	interactor := usecases.NewMemberInteractor(memberRepo, authUseCase)

	member := &domain.Member{
		MemberNumber: "M12345",
		Username:     "testuser",
		FullName:     "Test User",
		Email:        "testuser@example.com",
	}
	member.AssignPassword("password123")
	_ = memberRepo.Create(member)

	// When
	memberResponse, err := interactor.GetMyInfo("testuser")

	// Then
	assert.NoError(t, err)
	assert.Equal(t, member.FullName, memberResponse.FullName)
}

func TestMemberInteractor_GetMyInfo_Failure_NotFound(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	interactor := usecases.NewMemberInteractor(memberRepo, authUseCase)

	// When
	memberResponse, err := interactor.GetMyInfo("nonexistent")

	// Then
	assert.Error(t, err)
	assert.Nil(t, memberResponse)
}

func TestMemberInteractor_UpdateMyInfo_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	interactor := usecases.NewMemberInteractor(memberRepo, authUseCase)

	member := &domain.Member{
		MemberNumber: "M12345",
		Username:     "testuser",
		FullName:     "Old Name",
		Email:        "old@example.com",
	}
	member.AssignPassword("password123")
	_ = memberRepo.Create(member)

	updateReq := &request.UpdateMemberRequest{
		NickName: "New Name",
		Email:    "new@example.com",
	}

	// When
	err := interactor.UpdateMyInfo("testuser", updateReq)

	// Then
	assert.NoError(t, err)
	updatedMember, _ := memberRepo.GetByUserName("testuser")
	assert.Equal(t, "New Name", updatedMember.FullName)
	assert.Equal(t, "new@example.com", updatedMember.Email)
}

func TestMemberInteractor_UpdateMyInfo_Failure_UserNotFound(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	interactor := usecases.NewMemberInteractor(memberRepo, authUseCase)

	updateReq := &request.UpdateMemberRequest{
		NickName: "New Name",
		Email:    "new@example.com",
	}

	// When
	err := interactor.UpdateMyInfo("nonexistent", updateReq)

	// Then
	assert.Error(t, err)
}

func TestMemberInteractor_DeleteByUserName_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	interactor := usecases.NewMemberInteractor(memberRepo, authUseCase)

	member := &domain.Member{
		MemberNumber: "M12345",
		Username:     "testuser",
		FullName:     "Test User",
		Email:        "testuser@example.com",
	}
	member.AssignPassword("password123")
	_ = memberRepo.Create(member)

	// When
	err := interactor.DeleteByUserName("testuser")

	// Then
	assert.NoError(t, err)
	deletedMember, err := memberRepo.GetByUserName("testuser")
	assert.NoError(t, err)
	assert.True(t, deletedMember.IsWithdrawn)
}

func TestMemberInteractor_DeleteByUserName_Failure_UserNotFound(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	interactor := usecases.NewMemberInteractor(memberRepo, authUseCase)

	// When
	err := interactor.DeleteByUserName("nonexistent")

	// Then
	assert.Error(t, err)
}

func TestMemberInteractor_GetAllMembers_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	interactor := usecases.NewMemberInteractor(memberRepo, authUseCase)

	member1 := &domain.Member{
		MemberNumber: "M12345",
		Username:     "user1",
		FullName:     "User One",
		Email:        "user1@example.com",
	}
	member1.AssignPassword("password123")
	member2 := &domain.Member{
		MemberNumber: "M12346",
		Username:     "user2",
		FullName:     "User Two",
		Email:        "user2@example.com",
	}
	member2.AssignPassword("password123")
	_ = memberRepo.Create(member1)
	_ = memberRepo.Create(member2)

	// When
	members, err := interactor.GetAllMembers()

	// Then
	assert.NoError(t, err)
	assert.Len(t, members, 2)
}

func TestMemberInteractor_GetMemberStats_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	interactor := usecases.NewMemberInteractor(memberRepo, authUseCase)

	member1 := &domain.Member{
		ID:           12345,
		MemberNumber: "M12345",
		Username:     "user1",
		FullName:     "User One",
		Email:        "user1@example.com",
		IsWithdrawn:  false,
		CreatedAt:    time.Date(2024, 9, 10, 0, 0, 0, 0, time.UTC),
	}
	_ = member1.AssignPassword("password123")
	_ = memberRepo.Create(member1)

	member2 := &domain.Member{
		ID:           12346,
		MemberNumber: "M12346",
		Username:     "user2",
		FullName:     "User Two",
		Email:        "user2@example.com",
		IsWithdrawn:  false,
		CreatedAt:    time.Date(2024, 9, 15, 0, 0, 0, 0, time.UTC),
	}
	_ = member2.AssignPassword("password123")
	_ = memberRepo.Create(member2)

	// 멤버 삭제 수행
	_ = memberRepo.Delete(member2.ID)

	member2Fetched, _ := memberRepo.GetByID(member2.ID)
	withdrawnAt := time.Date(2024, 9, 20, 0, 0, 0, 0, time.UTC)
	member2Fetched.WithdrawnAt = &withdrawnAt
	_ = memberRepo.Update(member2Fetched)

	// When
	stats, err := interactor.GetMemberStats("2024-09")

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 2, stats["joined_members"])
	assert.Equal(t, 1, stats["deleted_members"])
}

func TestMemberInteractor_GetMemberStats_Failure_InvalidMonth(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	interactor := usecases.NewMemberInteractor(memberRepo, authUseCase)

	// When
	stats, err := interactor.GetMemberStats("invalid-month")

	// Then
	assert.Error(t, err)
	assert.Nil(t, stats)
}
