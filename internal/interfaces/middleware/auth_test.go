package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/HongJungWan/commerce-system/internal/interfaces/middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuthMiddleware_Success(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.JWTAuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		accountId := c.GetString("account_id")
		isAdmin := c.GetBool("is_admin")
		memberNumber := c.GetString("member_number")
		c.JSON(http.StatusOK, gin.H{
			"message":       "Access granted",
			"account_id":    accountId,
			"is_admin":      isAdmin,
			"member_number": memberNumber,
		})
	})

	// JWT 토큰 생성: 올바른 서명 키와 클레임 사용
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"account_id":    "testuser",
		"is_admin":      true,
		"member_number": "M12345",
		"exp":           time.Now().Add(time.Hour * 1).Unix(),
	})
	tokenString, err := token.SignedString([]byte("commerce-system"))
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusOK, resp.Code)
	var response map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Access granted", response["message"])
	assert.Equal(t, "testuser", response["account_id"])
	assert.Equal(t, true, response["is_admin"])
	assert.Equal(t, "M12345", response["member_number"])
}

func TestJWTAuthMiddleware_Failure_MissingToken(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.JWTAuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Access granted"})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Authorization header missing", response["error"])
}

func TestJWTAuthMiddleware_Failure_InvalidToken(t *testing.T) {
	// Given
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.JWTAuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Access granted"})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")

	// When
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid token", response["error"])
}
