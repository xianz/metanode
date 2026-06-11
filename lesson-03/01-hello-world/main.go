package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		time.Sleep(time.Second * 1)
		c.JSON(http.StatusOK, gin.H{
			"message": "asldkfjasdf",
		})
	})
	// 路径参数
	r.GET("/files/*allpath", func(c *gin.Context) {
		allpath := c.Param("allpath")
		c.JSON(http.StatusOK, gin.H{
			"path": allpath,
		})
	})

	// 查询参数
	r.GET("/query", func(c *gin.Context) {
		// url参数
		keyword := c.Query("keyword")
		size := c.DefaultQuery("size", "10")
		// 表单参数
		name := c.PostForm("name")
		email := c.PostForm("email")

		c.JSON(http.StatusOK, gin.H{
			"keyword": keyword,
			"size":    size,
			"name":    name,
			"email":   email,
		})
	})

	// json参数
	r.POST("/api/users", func(c *gin.Context) {
		type UserRequest struct {
			Name  string `json:"name" binding:"required"`
			Email string `json:"email" binding:"required,email"`
			Age   int    `json:"age" binding:"required,min=1,max=120"`
		}
		var userRequest UserRequest
		if err := c.ShouldBindJSON(&userRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"name":  userRequest.Name,
			"email": userRequest.Email,
			"age":   userRequest.Age,
		})
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/users", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "v1 users",
			})
		})
	}

	// 启动
	r.Run()
}
