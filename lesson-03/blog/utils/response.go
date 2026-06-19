package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Error   any    `json:"error"`
}

func Success(c *gin.Context, data any) {
	c.JSON(200, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string, err any) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Error:   err,
	})
}

// 单纯地返回错误响应
func ValidateError(c *gin.Context, err any) {
	c.JSON(http.StatusUnprocessableEntity, Response{
		Code:    http.StatusUnprocessableEntity,
		Message: "validate failed",
		Error:   err,
	})
}
