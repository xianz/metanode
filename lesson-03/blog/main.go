package main

import (
	"blog/config"
	"blog/handlers"
	"blog/middleware"
	"blog/models"
	"blog/services"
	"blog/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var cfg config.Config
var db *gorm.DB

func main() {

	gin.SetMode(cfg.Server.GinMode)

	userService := &services.UserService{Db: db}
	userHandler := &handlers.UserHandler{UserService: userService, JWTConfig: &cfg.Jwt}
	blogService := &services.BlogService{Db: db}
	blogHandler := &handlers.BlogHandler{BlogService: blogService}

	// 创建Gin实例
	r := gin.Default()
	// 全局中间件
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// 公开路由
	public := r.Group("/api/v1")
	{
		public.GET("/", showConfig)
		// 文章列表
		public.GET("/blog", blogHandler.ListArticles)
		// 文章详情
		public.GET("/blog/:id", blogHandler.GetArticle)

		//注册
		public.POST("/users", userHandler.Register)
		// 登录
		r.POST("/login", userHandler.Login)
		// 数据库迁移
		r.GET("/automigrate", automigrate)
	}

	// 需要认证的路由
	protected := r.Group("/api/v1")
	protected.Use(middleware.Auth([]byte(cfg.Jwt.Secret)))
	{
		// 创建文章
		r.PUT("/articles/:id", blogHandler.CreateArticle)
		// 更新文章
		r.POST("/articles/:id", blogHandler.UpdateArticle)
		// 删除文章
		r.DELETE("/articles/:id", blogHandler.DeleteArticle)
	}

	// // 退出
	// r.GET("/logout", logout)

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

// 数据库迁移
func automigrate(c *gin.Context) {
	err := db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	// 初始化数据
	defaultPassword := "123456"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("生成哈希失败！")
	}
	// 插入用户
	users := []*models.User{
		{Username: "zhangsan", Email: "aa@a.com", Password: string(hashedPassword)},
		{Username: "lisi", Email: "bbbbb@qq.com", Password: string(hashedPassword)},
	}
	result := db.Create(&users)
	if result.Error != nil {
		log.Fatal("创建用户失败！", result.Error)
	}
	log.Println("创建用户成功，影响记录数：", result.RowsAffected)

	// 插入文章
	posts := []*models.Post{
		{Title: "Post 1", Content: "Content 1", UserID: users[0].ID},
		{Title: "Post 2", Content: "Content 2", UserID: users[1].ID},
	}
	result = db.Create(&posts)
	if result.Error != nil {
		log.Fatal("创建文章失败！", result.Error)
	}
	log.Println("创建文章成功，影响记录数：", result.RowsAffected)

	// 插入评论
	comments := []*models.Comment{
		{Content: "Comment 1", PostID: posts[0].ID, UserID: users[0].ID},
		{Content: "Comment 2", PostID: posts[1].ID, UserID: users[1].ID},
	}
	result = db.Create(&comments)
	if result.Error != nil {
		log.Fatal("创建评论失败！", result.Error)
	}
	log.Println("创建评论成功，影响记录数：", result.RowsAffected)

	utils.Success(c, gin.H{
		"message": "Database migrated",
	})
}
