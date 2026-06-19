package utils

import (
	"blog/models"
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 数据库迁移
func Automigrate(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		err := db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
		if err != nil {
			log.Fatal("数据库迁移失败：", err)
			return
		}
		// 初始化数据
		defaultPassword := "123456"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("生成哈希失败：", err)
			return
		}
		// 插入用户
		users := []*models.User{
			{Username: "zhangsan", Email: "aa@a.com", Password: string(hashedPassword)},
			{Username: "lisi", Email: "bbbbb@qq.com", Password: string(hashedPassword)},
		}
		result := db.Create(&users)
		if result.Error != nil {
			log.Printf("创建用户失败：%v", result.Error)
		} else {
			log.Println("创建用户成功，影响记录数：", result.RowsAffected)
		}

		// 插入文章
		posts := []*models.Post{
			{Title: "Post 1", Content: "Content 1", UserID: users[0].ID},
			{Title: "Post 2", Content: "Content 2", UserID: users[1].ID},
		}
		result = db.Create(&posts)
		if result.Error != nil {
			log.Printf("创建文章失败：%v", result.Error)
		} else {
			log.Println("创建文章成功，影响记录数：", result.RowsAffected)
		}

		// 插入评论
		comments := []*models.Comment{
			{Content: "Comment 1", PostID: posts[0].ID, UserID: users[0].ID},
			{Content: "Comment 2", PostID: posts[1].ID, UserID: users[1].ID},
		}
		result = db.Create(&comments)
		if result.Error != nil {
			log.Printf("创建评论失败：%v", result.Error)
		} else {
			log.Println("创建评论成功，影响记录数：", result.RowsAffected)
		}

		Success(c, gin.H{
			"message": "Database migrated",
		})
	}
}
