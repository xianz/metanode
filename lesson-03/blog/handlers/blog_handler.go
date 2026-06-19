package handlers

import (
	"blog/models"
	"blog/services"
	"blog/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BlogHandler struct {
	BlogService *services.BlogService
}

func (bh *BlogHandler) CreateArticle(c *gin.Context) {
	var req models.Post
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, err)
		return
	}
	if req.Title == "" || req.Content == "" {
		utils.NewAppError(http.StatusBadRequest, "标题和内容不能为空")
		return
	}
	if _, err := bh.BlogService.CreateArticle(req); err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, "创建文章成功")
}

func (bh *BlogHandler) UpdateArticle(c *gin.Context) {

}

func (bh *BlogHandler) DeleteArticle(c *gin.Context) {
	_id := c.Param("id")
	req_id, err := strconv.Atoi(_id)
	if err != nil || req_id < 1 {
		utils.ValidateError(c, "文章ID格式错误")
		return
	}
	if err := bh.BlogService.DeleteArticle(req_id); err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, "删除文章成功")
}

func (bh *BlogHandler) GetArticle(c *gin.Context) {
	_id := c.Param("id")
	req_id, err := strconv.Atoi(_id)
	if err != nil || req_id < 1 {
		utils.ValidateError(c, "文章ID格式错误")
		return
	}
	article, err := bh.BlogService.GetArticle(req_id)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, article)
}

func (bh *BlogHandler) ListArticles(c *gin.Context) {
	pageSize := 5
	req_page := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(req_page)
	if err != nil {
		utils.ValidateError(c, "页码格式错误")
		return
	}

	rs, err := bh.BlogService.ListArticles(pageSize, page)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, rs)
}
