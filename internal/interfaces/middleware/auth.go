package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorization 헤더 확인
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// 토큰 파싱
		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 토큰 서명 방법 확인
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("commerce-system"), nil
		})

		// 토큰 유효성 검사
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 클레임 추출
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "유효하지 않은 토큰입니다."})
			c.Abort()
			return
		}

		// user_id 추출
		userID, ok := claims["user_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "유효하지 않은 토큰입니다."})
			c.Abort()
			return
		}

		// is_admin 추출
		isAdmin, _ := claims["is_admin"].(bool)

		// member_number 추출 (추가된 부분)
		memberNumber, ok := claims["member_number"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "유효하지 않은 토큰입니다."})
			c.Abort()
			return
		}

		// 컨텍스트에 값 설정
		c.Set("user_id", userID)
		c.Set("is_admin", isAdmin)
		c.Set("member_number", memberNumber) // 추가된 부분
		c.Next()
	}
}
