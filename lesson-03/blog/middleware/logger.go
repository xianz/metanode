package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// 继续业务流程
		c.Next()

		fmt.Printf("Request: [%s] %s %v\n", c.Request.Method, c.Request.URL.Path, time.Since(start))
	}
}
