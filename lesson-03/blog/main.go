package main

import (
	"blog/config"
	"blog/handlers"
	"blog/middleware"
	"blog/services"
	"blog/utils"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
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
		public.GET("/", healthCheck)
		// 文章列表
		public.GET("/blog", blogHandler.ListArticles)
		// 文章详情
		public.GET("/blog/:id", blogHandler.GetArticle)

		//注册
		public.POST("/users", userHandler.Register)
		// 登录
		r.POST("/login", userHandler.Login)
		// 数据库迁移
		r.GET("/automigrate", utils.Automigrate(db))
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
		log.Fatal("加载配置失败", err)
	}
	db = utils.GetSqliteDB(&cfg.Database.Sqlite)
	// log.Println("完成")
	// log.Printf("Configuration: %+v", cfg)
	// log.Printf("数据库：%+v\n", db)
	log.Printf("Load config cost: %v", time.Since(start))
}

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":    "ok",
		"config":    cfg,
		"timestamp": time.Now().Unix(),
	})
}
