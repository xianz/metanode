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

	//// 创建Gin实例
	route := gin.Default()
	//// 全局中间件
	route.Use(middleware.Logger())
	route.Use(middleware.CORS())

	//// 公开路由
	public := route.Group("/api/v1")
	public.GET("/", healthCheck)
	// 文章列表
	public.GET("/articles", blogHandler.ListArticles)
	// 文章详情
	public.GET("/articles/:id", blogHandler.GetArticle)

	//注册
	public.POST("/users", userHandler.Register)
	// 登录
	public.POST("/login", userHandler.Login)
	// 数据库迁移
	public.GET("/automigrate", utils.Automigrate(db))

	//// 需要认证的路由
	protected := route.Group("/api/v1")
	protected.Use(middleware.Auth([]byte(cfg.Jwt.Secret)))
	{
		// 创建文章
		protected.PUT("/articles", blogHandler.CreateArticle)
		// 更新文章
		protected.POST("/articles/:id", blogHandler.UpdateArticle)
		// 删除文章
		protected.DELETE("/articles/:id", blogHandler.DeleteArticle)
	}

	// // 退出
	// route.GET("/logout", logout)

	// 启动服务
	route.Run(fmt.Sprintf(":%d", cfg.Server.Port))
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
	log.Printf("加载配置完成，耗时：%v\n", time.Since(start))
}

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":    "ok",
		"config":    cfg,
		"timestamp": time.Now().Unix(),
	})
}
