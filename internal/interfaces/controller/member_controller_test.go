package controller_test

import (
	"bytes"
	"encoding/json"
	"github.com/HongJungWan/commerce-system/internal/domain"
	"github.com/HongJungWan/commerce-system/internal/infrastructure/repository"
	"github.com/HongJungWan/commerce-system/internal/interfaces/controller"
	"github.com/HongJungWan/commerce-system/internal/usecases"
	"github.com/HongJungWan/commerce-system/test/fixtures"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMemberController_Register_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	memberInteractor := usecases.NewMemberInteractor(memberRepo)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo) // JWT 시크릿 키 설정
	memberController := controller.NewMemberController(memberInteractor, authUseCase)

	router := gin.Default()
	router.POST("/register", memberController.Register)

	newMember := domain.Member{
		Username:       "newuser",
		HashedPassword: "password123",
		FullName:       "New User",
		Email:          "newuser@example.com",
	}

	requestBody, _ := json.Marshal(newMember)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusCreated, resp.Code)
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "회원 가입이 완료되었습니다.", response["message"])
	assert.NotEmpty(t, response["token"])
}

func TestMemberController_Register_Failure_InvalidRequest(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	memberInteractor := usecases.NewMemberInteractor(memberRepo)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	memberController := controller.NewMemberController(memberInteractor, authUseCase)

	router := gin.Default()
	router.POST("/register", memberController.Register)

	invalidRequest := map[string]interface{}{
		"user_id": 123, // 잘못된 타입의 값
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
	memberInteractor := usecases.NewMemberInteractor(memberRepo)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	memberController := controller.NewMemberController(memberInteractor, authUseCase)

	member := &domain.Member{
		Username:       "testuser",
		HashedPassword: "password123",
		FullName:       "Test User",
		Email:          "testuser@example.com",
		MemberNumber:   "M12345",
	}
	_ = memberRepo.Create(member)

	router := gin.Default()
	router.GET("/me", func(c *gin.Context) {
		c.Set("user_id", "testuser")
		memberController.GetMyInfo(c)
	})

	req, _ := http.NewRequest("GET", "/me", nil)

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusOK, resp.Code)
	var retrievedMember domain.Member
	err := json.Unmarshal(resp.Body.Bytes(), &retrievedMember)
	assert.NoError(t, err)
	assert.Equal(t, "Test User", retrievedMember.FullName)
}

func TestMemberController_GetMyInfo_Failure_UserNotFound(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	memberInteractor := usecases.NewMemberInteractor(memberRepo)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	memberController := controller.NewMemberController(memberInteractor, authUseCase)

	router := gin.Default()
	router.GET("/me", func(c *gin.Context) {
		c.Set("user_id", "nonexistent")
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
	memberInteractor := usecases.NewMemberInteractor(memberRepo)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	memberController := controller.NewMemberController(memberInteractor, authUseCase)

	member := &domain.Member{
		Username:       "testuser",
		HashedPassword: "password123",
		FullName:       "Old Name",
		Email:          "old@example.com",
	}
	_ = memberRepo.Create(member)

	router := gin.Default()
	router.PUT("/me", func(c *gin.Context) {
		c.Set("user_id", "testuser")
		memberController.UpdateMyInfo(c)
	})

	updateData := map[string]interface{}{
		"name":  "New Name",
		"email": "new@example.com",
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

	updatedMember, _ := memberRepo.GetByUserName("testuser")
	assert.Equal(t, "New Name", updatedMember.FullName)
	assert.Equal(t, "new@example.com", updatedMember.Email)
}

func TestMemberController_UpdateMyInfo_Failure_InvalidRequest(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	memberInteractor := usecases.NewMemberInteractor(memberRepo)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	memberController := controller.NewMemberController(memberInteractor, authUseCase)

	router := gin.Default()
	router.PUT("/me", func(c *gin.Context) {
		c.Set("user_id", "testuser")
		memberController.UpdateMyInfo(c)
	})

	invalidRequest := map[string]interface{}{
		"name": 123, // 잘못된 타입의 값
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
	memberInteractor := usecases.NewMemberInteractor(memberRepo)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	memberController := controller.NewMemberController(memberInteractor, authUseCase)

	member := &domain.Member{
		Username:       "testuser",
		HashedPassword: "password123",
		FullName:       "Test User",
		Email:          "testuser@example.com",
	}
	_ = memberRepo.Create(member)

	router := gin.Default()
	router.DELETE("/me", func(c *gin.Context) {
		c.Set("user_id", "testuser")
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

	deletedMember, err := memberRepo.GetByUserName("testuser")
	assert.Error(t, err)
	assert.Nil(t, deletedMember)
}

func TestMemberController_GetAllMembers_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	memberInteractor := usecases.NewMemberInteractor(memberRepo)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	memberController := controller.NewMemberController(memberInteractor, authUseCase)

	member1 := &domain.Member{
		Username:       "user1",
		HashedPassword: "password123",
		FullName:       "User One",
		Email:          "user1@example.com",
	}
	member2 := &domain.Member{
		Username:       "user2",
		HashedPassword: "password123",
		FullName:       "User Two",
		Email:          "user2@example.com",
	}
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
	var members []domain.Member
	err := json.Unmarshal(resp.Body.Bytes(), &members)
	assert.NoError(t, err)
	assert.Len(t, members, 1)
}

func TestMemberController_GetAllMembers_Failure_Unauthorized(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	memberInteractor := usecases.NewMemberInteractor(memberRepo)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
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
