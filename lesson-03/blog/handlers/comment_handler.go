package handlers

import (
	"blog/services"
	"blog/utils"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	CommentService *services.CommentService
}

func (ch *CommentHandler) CreateComment(c *gin.Context) {
	// var req models.CommentRequest
	var req map[string]any
	// 得到content
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, err)
		return
	}
	// 得到post_id
	post_id_str := c.Query("post_id")
	if post_id_str == "" {
		utils.ValidateError(c, "文章编号为空")
		return
	}
	post_id, err := strconv.ParseInt(post_id_str, 10, 64)
	if err != nil {
		utils.ValidateError(c, "文章编号有误")
		return
	}
	// 使用当前登录用户
	userID, exists := c.Get("userID")
	if !exists {
		utils.HandleError(c, errors.New("怎么回事？"))
		return
	}
	userIDUint, ok := userID.(uint)
	if !ok {
		utils.HandleError(c, errors.New("userID类型错误"))
		return
	}
	// 组装comment
	comment := make(map[string]any)
	comment["content"] = req["content"]
	comment["PostID"] = post_id
	comment["UserID"] = userIDUint
	rowsAffected, err := ch.CommentService.CreateComment(comment)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, gin.H{
		"rows_affected": rowsAffected,
	})
}

func (ch *CommentHandler) GetComments(c *gin.Context) {
	req_postID := c.Query("post_id")
	if req_postID == "" {
		utils.ValidateError(c, "文章编号为空")
		return
	}
	post_id, err := strconv.ParseInt(req_postID, 10, 64)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	// log.Fatalf("%T : %+v\n", post_id, post_id)
	responseData, err := ch.CommentService.GetCommentsByPostID(post_id)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, responseData)
}
