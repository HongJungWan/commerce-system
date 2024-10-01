package controller

import (
	"net/http"

	"github.com/HongJungWan/commerce-system/internal/interfaces/dto/request"
	"github.com/HongJungWan/commerce-system/internal/interfaces/dto/response"
	"github.com/HongJungWan/commerce-system/internal/usecases"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authUseCase *usecases.AuthUseCase
}

func NewAuthController(authUseCase *usecases.AuthUseCase) *AuthController {
	return &AuthController{authUseCase: authUseCase}
}

// Login godoc
// @Summary      사용자 로그인
// @Description  사용자 인증 정보를 확인하고 JWT 토큰을 반환합니다.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        loginRequest body request.LoginRequest true "사용자 로그인 정보"
// @Success      200 {object} response.LoginResponse "로그인 성공"
// @Failure      400 {object} map[string]string "잘못된 요청"
// @Failure      401 {object} map[string]string "인증 실패"
// @Failure      500 {object} map[string]string "서버 오류"
// @Router       /login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var loginRequest request.LoginRequest

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

	loginResponse := response.LoginResponse{
		Token: token,
	}

	c.JSON(http.StatusOK, loginResponse)
}
