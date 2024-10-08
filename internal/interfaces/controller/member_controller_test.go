package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/infrastructure/repository"
	"github.com/HongJungWan/commerce-system/internal/interfaces/controller"
	"github.com/HongJungWan/commerce-system/internal/interfaces/dto/request"
	"github.com/HongJungWan/commerce-system/internal/interfaces/dto/response"
	"github.com/HongJungWan/commerce-system/internal/usecases"
	"github.com/HongJungWan/commerce-system/test/fixtures"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMemberController_Register_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo) // JWT 시크릿 키 설정
	memberInteractor := usecases.NewMemberInteractor(memberRepo, authUseCase)
	memberController := controller.NewMemberController(memberInteractor, authUseCase)

	router := gin.Default()
	router.POST("/register", memberController.Register)

	newMember := request.CreateMemberRequest{
		AccountId: "newuser",
		Password:  "password123",
		NickName:  "New User",
		Email:     "newuser@example.com",
	}

	requestBody, _ := json.Marshal(newMember)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusCreated, resp.Code)
	var responseData response.RegisterMemberResponse
	err := json.Unmarshal(resp.Body.Bytes(), &responseData)
	assert.NoError(t, err)
	assert.Equal(t, "회원 가입이 완료되었습니다.", responseData.Message)
	assert.Equal(t, "newuser", responseData.User.Username)
}

func TestMemberController_Register_Failure_InvalidRequest(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	memberInteractor := usecases.NewMemberInteractor(memberRepo, authUseCase)
	memberController := controller.NewMemberController(memberInteractor, authUseCase)

	router := gin.Default()
	router.POST("/register", memberController.Register)

	invalidRequest := map[string]interface{}{
		"account_id": 0, // 잘못된 타입의 값
	}

	requestBody, _ := json.Marshal(invalidRequest)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestMemberController_GetMyInfo_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	memberInteractor := usecases.NewMemberInteractor(memberRepo, authUseCase)
	memberController := controller.NewMemberController(memberInteractor, authUseCase)

	member := &domain.Member{
		MemberNumber: "M1234",
		AccountId:    "testuser",
		NickName:     "Test User",
		Email:        "testuser@example.com",
	}
	member.AssignPassword("password123")
	_ = memberRepo.Create(member)

	router := gin.Default()
	router.GET("/me", func(c *gin.Context) {
		c.Set("account_id", "testuser")
		memberController.GetMyInfo(c)
	})

	req, _ := http.NewRequest("GET", "/me", nil)

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusOK, resp.Code)
	var responseData response.MemberResponse
	err := json.Unmarshal(resp.Body.Bytes(), &responseData)
	assert.NoError(t, err)
	assert.Equal(t, "Test User", responseData.FullName)
}

func TestMemberController_GetMyInfo_Failure_UserNotFound(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	memberInteractor := usecases.NewMemberInteractor(memberRepo, authUseCase)
	memberController := controller.NewMemberController(memberInteractor, authUseCase)

	router := gin.Default()
	router.GET("/me", func(c *gin.Context) {
		c.Set("username", "nonexistent")
		memberController.GetMyInfo(c)
	})

	req, _ := http.NewRequest("GET", "/me", nil)

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestMemberController_UpdateMyInfo_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	memberInteractor := usecases.NewMemberInteractor(memberRepo, authUseCase)
	memberController := controller.NewMemberController(memberInteractor, authUseCase)

	member := &domain.Member{
		MemberNumber: "M1234",
		AccountId:    "testuser",
		NickName:     "Old Name",
		Email:        "old@example.com",
	}
	member.AssignPassword("password123")
	_ = memberRepo.Create(member)

	router := gin.Default()
	router.PUT("/me", func(c *gin.Context) {
		c.Set("account_id", "testuser")
		memberController.UpdateMyInfo(c)
	})

	updateData := request.UpdateMemberRequest{
		NickName: "New Name",
		Email:    "new@example.com",
	}
	requestBody, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", "/me", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusOK, resp.Code)
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "정보가 수정되었습니다.", response["message"])

	updatedMember, _ := memberRepo.GetByAccountId("testuser")
	assert.Equal(t, "New Name", updatedMember.NickName)
	assert.Equal(t, "new@example.com", updatedMember.Email)
}

func TestMemberController_UpdateMyInfo_Failure_InvalidRequest(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	memberInteractor := usecases.NewMemberInteractor(memberRepo, authUseCase)
	memberController := controller.NewMemberController(memberInteractor, authUseCase)

	member := &domain.Member{
		MemberNumber: "M1234",
		AccountId:    "testuser",
		NickName:     "Test User",
		Email:        "testuser@example.com",
	}
	member.AssignPassword("password123")
	memberRepo.Create(member)

	router := gin.Default()
	router.PUT("/me", func(c *gin.Context) {
		c.Set("username", "testuser")
		memberController.UpdateMyInfo(c)
	})

	invalidRequest := map[string]interface{}{
		"nick_name": 123, // 잘못된 타입
	}
	requestBody, _ := json.Marshal(invalidRequest)
	req, _ := http.NewRequest("PUT", "/me", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestMemberController_DeleteMyAccount_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	memberInteractor := usecases.NewMemberInteractor(memberRepo, authUseCase)
	memberController := controller.NewMemberController(memberInteractor, authUseCase)

	member := &domain.Member{
		MemberNumber: "M1234",
		AccountId:    "testuser",
		NickName:     "Test User",
		Email:        "testuser@example.com",
	}
	member.AssignPassword("password123")
	_ = memberRepo.Create(member)

	router := gin.Default()
	router.DELETE("/me", func(c *gin.Context) {
		c.Set("account_id", "testuser")
		memberController.DeleteMyAccount(c)
	})

	req, _ := http.NewRequest("DELETE", "/me", nil)

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusOK, resp.Code)
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "계정이 삭제되었습니다.", response["message"])

	deletedMember, err := memberRepo.GetByAccountId("testuser")
	assert.NoError(t, err)
	assert.True(t, deletedMember.IsWithdrawn)
}

func TestMemberController_GetAllMembers_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	memberInteractor := usecases.NewMemberInteractor(memberRepo, authUseCase)
	memberController := controller.NewMemberController(memberInteractor, authUseCase)

	member1 := &domain.Member{
		MemberNumber: "M1234",
		AccountId:    "user1",
		NickName:     "User One",
		Email:        "user1@example.com",
	}
	member1.AssignPassword("password123")
	member2 := &domain.Member{
		MemberNumber: "M1235",
		AccountId:    "user2",
		NickName:     "User Two",
		Email:        "user2@example.com",
	}
	member2.AssignPassword("password123")
	_ = memberRepo.Create(member1)
	_ = memberRepo.Create(member2)

	router := gin.Default()
	router.GET("/members", func(c *gin.Context) {
		c.Set("is_admin", true)
		memberController.GetAllMembers(c)
	})

	req, _ := http.NewRequest("GET", "/members", nil)

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusOK, resp.Code)
	var members []*response.MemberResponse
	err := json.Unmarshal(resp.Body.Bytes(), &members)
	assert.NoError(t, err)
	assert.Len(t, members, 2)
}

func TestMemberController_GetAllMembers_Failure_Unauthorized(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	memberInteractor := usecases.NewMemberInteractor(memberRepo, authUseCase)
	memberController := controller.NewMemberController(memberInteractor, authUseCase)

	router := gin.Default()
	router.GET("/members", func(c *gin.Context) {
		c.Set("is_admin", false)
		memberController.GetAllMembers(c)
	})

	req, _ := http.NewRequest("GET", "/members", nil)

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusForbidden, resp.Code)
}
