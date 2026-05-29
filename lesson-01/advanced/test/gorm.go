package main

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:100"`
	Age  int
}

func main() {
	// 使用纯 Go 驱动，无需 CGO
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败: " + err.Error())
	}

	// 自动迁移
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("迁移失败: " + err.Error())
	}

	// 测试操作
	db.Create(&User{Name: "测试用户", Age: 18})

	var user User
	db.First(&user, 1)
	fmt.Printf("用户: %+v\n", user)

	fmt.Println("数据库操作成功！")
}
