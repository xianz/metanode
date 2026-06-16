package main

import (
	"blog/config"
	"blog/middleware"
	"blog/models"
	"blog/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var cfg config.Config
var db *gorm.DB

func main() {
	gin.SetMode(cfg.Server.GinMode)

	r := gin.Default()
	// 默认路由
	r.GET("/", showConfig)

	// 博客相关路由
	r.GET("/blog/list", blogList)
	r.GET("/blog/detail", blogDetail)

	// 数据库迁移
	r.GET("/automigrate", automigrate)

	// 登录
	r.POST("/login", login)

	// // 注册
	// r.GET("/register", register)

	// // 退出
	// r.GET("/logout", logout)

	// // 博客管理
	g1 := r.Group("api")
	g1.Use(middleware.AuthMiddleware())
	// 博客编辑
	g1.GET("/blog/edit", blogEdit)
	// // 博客删除
	// g1.GET("/blog/delete", blogDelete)
	// // 博客创建
	// g1.GET("/blog/create", blogCreate)
	// // 博客更新
	// g1.GET("/blog/update", blogUpdate)

	// 启动服务
	r.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}

// 初始化配置和数据库
func init() {
	// log.Println("初始工作...")
	start := time.Now()
	err := config.LoadConfig(&cfg)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	db = utils.GetSqliteDB(&cfg.Database.Sqlite)
	// log.Println("完成")
	// log.Printf("Configuration: %+v", cfg)
	log.Printf("数据库：%+v\n", db)
	log.Printf("Load config cost: %v", time.Since(start))
}

func showConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"config": cfg,
	})
}

// 获取博客列表
func blogList(c *gin.Context) {
	pageSize := 5
	pageStr, _ := c.GetQuery("page")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	var posts []models.Post
	scopePage(page, pageSize).Preload("Comments").Find(&posts)
	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}

// 获取博客详情
func blogDetail(c *gin.Context) {
	postIDStr, _ := c.GetQuery("id")
	postID, _ := strconv.Atoi(postIDStr)
	if postID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid post ID",
		})
		return
	}
	var post models.Post
	err := db.Preload("Comments").First(&post, postID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Post not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

// 分页查询
func scopePage(page int, pageSize int) *gorm.DB {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 5
	}
	return db.Offset((page - 1) * pageSize).Limit(pageSize)
}

func blogEdit(c *gin.Context) {
	reqPostID := c.PostForm("postID")
	postID, err := strconv.Atoi(reqPostID)
	postContent := c.PostForm("postContent")
	if err != nil || postID < 1 || len(postContent) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "青检查请求参数",
		})
		return
	}
	result := db.Model(models.Post{}).Where("id=?", postID).Update("content", postContent)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": result.Error,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":     "success",
		"rowAffected": result.RowsAffected,
	})
	// db.Model(&models.Post{}).Where("id = ?", reqPostID).Update("content", reqPostContent)

}

// 数据库迁移
func automigrate(c *gin.Context) {
	// if utils.SqliteFileExists(cfg.Database.Sqlite.Path) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "Database already exists",
	// 	})
	// 	return
	// }
	err := db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}, &models.Tag{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	// 初始化数据
	users := []*models.User{
		{Username: "zhangsan", Password: "123456"},
		{Username: "lisi", Password: "123456"},
	}
	db.Create(&users)

	posts := []*models.Post{
		{Title: "Post 1", Content: "Content 1", UserID: users[0].ID},
		{Title: "Post 2", Content: "Content 2", UserID: users[1].ID},
	}
	db.Create(&posts)
	comments := []*models.Comment{
		{Content: "Comment 1", PostID: posts[0].ID, UserID: users[0].ID},
		{Content: "Comment 2", PostID: posts[1].ID, UserID: users[1].ID},
	}
	db.Create(&comments)

	c.JSON(http.StatusOK, gin.H{
		"message": "Database migrated",
		"data": gin.H{
			"user.rowAffected": db.RowsAffected,
		},
	})
}

// 登录
func login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Username and password are required",
		})
		return
	}
	var user models.User
	db.Where("username=? AND password=?", username, password).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid username or password",
		})
		return
	}
	// 生成jwt
	token, err := (&utils.JWTManager{}).GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  gin.H{"userID": user.ID, "username": user.Username},
	})
}
