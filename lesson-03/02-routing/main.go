package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	/////////////// 查询参数绑定
	type QueryParams struct {
		Name string `form:"name" binding:"required"`
		Age  int    `form:"age" binding:"required,gte=1,lte=120"`
		Sex  string `form:"sex" binding:"eq=male|eq=female"`
	}
	queryParams := QueryParams{}
	r.GET("/api/products", func(c *gin.Context) {
		if err := c.ShouldBindQuery(&queryParams); err != nil {
			errResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		success(c, http.StatusOK, queryParams)
	})

	//////////////////表单验证
	r.POST("/api/login", func(c *gin.Context) {
		userRequest := &UserRequest{}
		if req := c.ShouldBind(userRequest); req != nil {
			errResponse(c, http.StatusBadRequest, req.Error())
			return
		} else if userRequest.Username != "admin" || userRequest.Password != "pwd" {
			errResponse(c, http.StatusUnauthorized, "账号密码错误")
			return
		}
		success(c, http.StatusOK, userRequest)
	})

	/////////////////// 重定向
	r.GET("/redirect", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "https://www.baidu.com")
	})

	/////////////////// 数据流
	r.GET("/data", func(ctx *gin.Context) {
		ctx.Data(http.StatusOK, "application/octet-stream", []byte("hello world"))
	})

	/////////////////// 分组
	g1 := r.Group("/api")
	{
		g1.GET("/user", func(ctx *gin.Context) {
			success(ctx, http.StatusOK, "user")
		})
	}

	r.Run()
}

type UserRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func success(c *gin.Context, httpStatus int, data interface{}) {
	c.JSON(httpStatus, gin.H{
		"status":  200,
		"message": "success",
		"data":    data,
	})
}

func errResponse(c *gin.Context, httpStatus int, message string) {
	c.JSON(httpStatus, gin.H{
		"status":  httpStatus,
		"message": message,
	})
}
