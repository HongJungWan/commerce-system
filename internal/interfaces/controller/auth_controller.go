package controller

import (
	"github.com/HongJungWan/commerce-system/internal/usecases"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthController struct {
	authUseCase *usecases.AuthUseCase
}

func NewAuthController(authUseCase *usecases.AuthUseCase) *AuthController {
	return &AuthController{authUseCase: authUseCase}
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var loginRequest struct {
		UserID   string `json:"user_id"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// FIXME: 인증 로직 추가 (예: DB에서 사용자 조회 등)
	// FIXME: ...

	token, err := ctrl.authUseCase.GenerateToken(loginRequest.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
