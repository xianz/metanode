package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.Use(loggerMiddleware())
	r.Use(gin.Recovery())
	// r.Use(recoveryMiddlerware())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})

	// 报个错
	r.GET("/panic", func(ctx *gin.Context) {
		var arr []string
		arr[1] = "test"
	})

	// 分个组
	g1 := r.Group("/v1")
	g1.Use(authMiddleware())
	{
		g1.GET("/login", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			c.JSON(200, gin.H{
				"userID": userID,
			})
		})
	}

	// 子分组 + 测试顺序
	g1_1 := r.Group("/v1/point")
	g1_1.Use(middleware1(), middleware2(), middleware3())
	{
		g1_1.GET("/", func(c *gin.Context) {
			c.String(200, "handle执行结束")
			fmt.Println("handle 执行了")
		})
	}

	// 限制连接数
	g2 := r.Group("/v2")
	g2.Use(rateLimitMiddleware())
	{
		g2.GET("/limit", func(c *gin.Context) {
			c.String(200, "limit test")
		})
	}

	r.Run()
}

// ========== 断点 顺序 ==========
func middleware1() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("断点1开始")
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		fmt.Printf("断点1结束，耗时: %v\n", latency)
	}
}
func middleware2() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("断点2开始")
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		fmt.Printf("断点2结束，耗时: %v\n", latency)
	}
}
func middleware3() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("断点3开始")
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		fmt.Printf("断点3结束，耗时: %v\n", latency)
	}
}

func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		fmt.Printf("【%s】%s %d %v\n", method, path, status, latency)
	}
}

func recoveryMiddlerware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		if recovered != nil {
			fmt.Println("发生错误了")
		}
		c.String(http.StatusInternalServerError, "internal server error")
		c.Abort()

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

// ========== 限流中间件（简单示例） ==========
var requestCount = make(map[string]int)
var lastReset = time.Now()

func rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		// 每分钟重置一次
		if now.Sub(lastReset) > time.Minute {
			requestCount = make(map[string]int)
			lastReset = now
		}

		// 检查请求次数
		if requestCount[ip] >= 5 {
			fmt.Println("too many requests from IP:", ip)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			c.Abort()
			return
		}
		requestCount[ip]++
		fmt.Printf("IP: %s, 请求次数: %d\n", ip, requestCount[ip])

		c.Next()
	}
}
