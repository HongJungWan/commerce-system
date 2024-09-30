package controller

import (
	"net/http"

	"github.com/HongJungWan/commerce-system/internal/usecases"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authUseCase *usecases.AuthUseCase
}

func NewAuthController(authUseCase *usecases.AuthUseCase) *AuthController {
	return &AuthController{authUseCase: authUseCase}
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var loginRequest struct {
		Username       string `json:"username"`
		HashedPassword string `json:"hashed_password"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	member, err := ctrl.authUseCase.Authenticate(loginRequest.Username, loginRequest.HashedPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := ctrl.authUseCase.GenerateToken(member)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
