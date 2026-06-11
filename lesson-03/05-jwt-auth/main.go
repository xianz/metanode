package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte("my-secret-key")

func generateToken(userID int, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
	// claims := jwt.MapClaims{
	// 	"user_id":  userID,
	// 	"username": username,
	// 	"exp":      time.Now().Add(time.Hour * 24).Unix(),
	// }
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// return token.SignedString(jwtSecret)
}

// func generateToken(username string) (string, error) {
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
// 		Username: username,
// 		Password: "password",
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
// 			IssuedAt:  jwt.NewNumericDate(time.Now()),
// 		},
// 	})
// 	return token.SignedString(jwtSecret)
// }

func main() {
	fmt.Println("JWT Auth")
	r := gin.Default()
	r.GET("/login", func(c *gin.Context) {
		c.String(200, "Login page")
	})
	v1 := r.Group("v1")
	v1.Use(authMiddleware())
	v1.GET("/protected", func(c *gin.Context) {
		c.String(200, "Protected page")
	})
}

// ========== 认证中间件（示例） ==========
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// 模拟验证 Token
		if token != "Bearer valid-token" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到 Context
		c.Set("userID", 1)
		c.Set("username", "admin")

		c.Next()
	}
}
