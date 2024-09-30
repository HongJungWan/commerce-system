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
	"github.com/HongJungWan/commerce-system/internal/usecases"
	"github.com/HongJungWan/commerce-system/test/fixtures"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthController_Login_Success(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	authController := controller.NewAuthController(authUseCase)

	member := &domain.Member{
		Username:       "testuser",
		HashedPassword: "password123",
		FullName:       "Test User",
		Email:          "testuser@example.com",
	}
	_ = member.SetPassword(member.HashedPassword)
	_ = memberRepo.Create(member)

	router := gin.Default()
	router.POST("/login", authController.Login)

	loginRequest := map[string]string{
		"username":        "testuser",
		"hashed_password": "password123",
	}
	requestBody, _ := json.Marshal(loginRequest)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusOK, resp.Code)
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["token"])
}

func TestAuthController_Login_Failure_InvalidCredentials(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	authController := controller.NewAuthController(authUseCase)

	router := gin.Default()
	router.POST("/login", authController.Login)

	loginRequest := map[string]string{
		"username":        "nonexistent",
		"hashed_password": "password123",
	}
	requestBody, _ := json.Marshal(loginRequest)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid credentials", response["error"])
}

func TestAuthController_Login_Failure_InvalidRequest(t *testing.T) {
	// Given
	db := fixtures.SetupTestDB()
	memberRepo := repository.NewMemberRepository(db)
	authUseCase := usecases.NewAuthUseCase("secret", memberRepo)
	authController := controller.NewAuthController(authUseCase)

	router := gin.Default()
	router.POST("/login", authController.Login)

	invalidRequest := map[string]interface{}{
		"username": 123, // 잘못된 타입의 값
	}
	requestBody, _ := json.Marshal(invalidRequest)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}
