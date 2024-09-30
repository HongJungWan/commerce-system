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

		// Username 추출
		Username, ok := claims["username"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "유효하지 않은 토큰입니다."})
			c.Abort()
			return
		}

		// is_admin 추출 (타입 변환 수정)
		isAdmin := false
		if val, ok := claims["is_admin"]; ok {
			fmt.Printf("is_admin type: %T, value: %v\n", val, val)
			switch v := val.(type) {
			case bool:
				isAdmin = v
			case float64:
				isAdmin = v != 0
			case string:
				isAdmin = v == "true"
			default:
				isAdmin = false
			}
		}

		// member_number 추출
		memberNumber, ok := claims["member_number"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "유효하지 않은 토큰입니다."})
			c.Abort()
			return
		}

		// 컨텍스트에 값 설정
		c.Set("username", Username)
		c.Set("is_admin", isAdmin)
		c.Set("member_number", memberNumber)
		c.Next()
	}
}
