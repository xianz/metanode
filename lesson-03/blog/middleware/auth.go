package middleware

import (
	"blog/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth(jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		headerAuth := c.GetHeader("Authorization")
		//// 不能为空
		if headerAuth == "" {
			utils.Error(c, http.StatusUnauthorized, "Authorization header required", nil)
			c.Abort()
			return
		}
		//// 开始验证解码
		// 提取
		parts := strings.Split(headerAuth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Error(c, http.StatusUnauthorized, "Invalid Authorization header format", nil)
			c.Abort()
			return
		}
		// 验证解码
		token := parts[1]
		claims, err := utils.ParseToken(jwtSecret, token)
		if err != nil {
			utils.Error(c, http.StatusUnauthorized, "Invalid token", err)
			c.Abort()
			return
		}
		// 将用户信息存储到 Context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		// 继续处理请求
		c.Next()
	}
}
