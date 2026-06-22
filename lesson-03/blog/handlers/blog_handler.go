package handlers

import (
	"blog/models"
	"blog/services"
	"blog/utils"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BlogHandler struct {
	BlogService *services.BlogService
}

type ParamsRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// 测试 id
func (bh *BlogHandler) CreateArticle(c *gin.Context) {
	var req models.PostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, err)
		return
	}
	if req.Title == "" || req.Content == "" {
		utils.ValidateError(c, "标题和内容不能为空")
		return
	}
	userID, exists := c.Get("userID")
	if !exists {
		utils.HandleError(c, fmt.Errorf("请检查登录是否失效"))
		// utils.ValidateError("。。。") // 这个错误不应该出现，出现了就是程序有bug了，所以直接丢给后端好了不用反应给前端
		return
	}
	// log.Fatalf("userID类型：%T 值：%+v\n", userID, userID)
	_userID, ok := userID.(uint)
	if !ok {
		utils.HandleError(c, fmt.Errorf("用户ID类型错误，请检查中间件赋值是否正确"))
		return
	}
	req.UserID = _userID
	// log.Fatalf("%+v\n", req)
	// log.Printf("@@@@@ %+v\n", req)
	if rowAffected, err := bh.BlogService.CreateArticle(req); err != nil {
		utils.HandleError(c, err)
		return
	} else {
		utils.Success(c, gin.H{
			"row_affected": rowAffected,
		})
	}
}

// 更新post done
func (bh *BlogHandler) UpdateArticle(c *gin.Context) {
	var reqParams ParamsRequest
	if err := c.ShouldBindUri(&reqParams); err != nil {
		utils.ValidateError(c, "传递参数有误，请检查")
		return
	}
	var req models.PostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, err)
		return
	}

	// utils.LogWithLocation(fmt.Sprintf("post_id: %T, %v", post_id, post_id))
	// log.Fatal()
	// log.Fatalf("-------%T , %+v\n", post_id, post_id)

	if req.Title == "" && req.Content == "" {
		utils.ValidateError(c, "标题和内容不能同时为空，或更新ID不能为空")
		return
	}

	rowAffected, err := bh.BlogService.UpdateArticle(reqParams.ID, req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, gin.H{
		"row_affected": rowAffected,
	})
}

func (bh *BlogHandler) DeleteArticle(c *gin.Context) {
	var reqParams ParamsRequest
	if err := c.ShouldBindUri(&reqParams); err != nil {
		utils.ValidateError(c, "传递参数有误，请检查")
		return
	}
	rowAffected, err := bh.BlogService.DeleteArticle(reqParams.ID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, gin.H{
		"row_affected": rowAffected,
	})
}

func (bh *BlogHandler) GetArticle(c *gin.Context) {
	var reqParams ParamsRequest
	if err := c.ShouldBindUri(&reqParams); err != nil {
		utils.ValidateError(c, "传递参数有误，请检查")
		return
	}
	article, err := bh.BlogService.GetArticle(reqParams.ID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, article)
}

func (bh *BlogHandler) ListArticles(c *gin.Context) {
	reqPage := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(reqPage)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	reqSize := c.DefaultQuery("pageSize", "5")
	pageSize, err := strconv.Atoi(reqSize)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 5
	}

	rs, err := bh.BlogService.ListArticles(pageSize, page)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, rs)
}
