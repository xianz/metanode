package handlers

import (
	"blog/config"
	"blog/models"
	"blog/services"
	"blog/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService *services.UserService
	JWTConfig   *config.JwtConfig
}

func (h *UserHandler) Login(c *gin.Context) {
	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil {
		// fmt.Printf("%#v\n", err.Error())
		utils.HandleError(c, err)
		// utils.ValidateError(c, parseValidationErrors(err))
		return
	}

	// 验证用户登录
	if req.Username == "" || req.Password == "" {
		utils.ValidateError(c, "用户名或密码不能为空")
		return
	}

	// 验证用户登录
	user, err := h.UserService.Authenticate(req.Username, req.Password)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// 登录成功
	utils.Success(c, map[string]any{
		"user": user,
	})
}

func (h *UserHandler) Register(c *gin.Context) {
	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, err)
		return
	}

	// 验证用户注册
	if req.Username == "" || req.Password == "" {
		utils.ValidateError(c, "用户名或密码不能为空")
		return
	}

	// 验证用户注册
	user, err := h.UserService.CreateUser(req.Username, req.Password)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// 注册成功
	utils.Success(c, map[string]any{
		"user": user,
	})
}

func parseValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	// 简化处理，实际应该解析 binding 错误
	errors["general"] = err.Error()
	return errors
}
