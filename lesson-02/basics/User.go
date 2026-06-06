package main

import (
	// "lesson02/basics/utils"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:64;not null"`
	Email     string    `gorm:"size:128;uniqueIndex;not null"`
	Age       uint8     `gorm:"not null"`
	Status    string    `gorm:"size:16;default:active;index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func SetupDB(dbName string) (*gorm.DB, error) {
	dbPath, err := GenerateSQLiteFilePath("users.db")
	if err != nil {
		log.Fatalln(err)
	}
	// log.Fatalln(dbPath)
	_, err = os.Stat(dbPath)
	dbExists := os.IsNotExist(err)
	db := OpenSqlite(dbPath)
	if dbExists {
		// 不存在就建表（自动迁移）
		fmt.Println(dbPath, "@@@数据库文件不存在，执行自动迁移，和数据插入")
		migrateErr := db.AutoMigrate(&User{})
		if migrateErr != nil {
			// log.Fatalf("auto migrate: %v", err)
			return nil, err
		}
		if err := db.Create(users).Error; err != nil {
			return nil, err
		}
	}
	return db, nil
}

// 新加字段
func UpdateDB(db *gorm.DB) error {
	type User struct {
		Phone       string    `gorm:"uniqueIndex"`
		LastLoginAt time.Time `gorm:"autoCreateTime;index"`
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		return err
	}
	return nil
}

// // 创建账户
func CreateUser(db *gorm.DB, name, email string) (*User, error) {
	user := &User{
		Name: name, Email: email, Status: "active",
	}
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// // 查找emial用户
func SearchUsersByEmail(db *gorm.DB, emailPattern string, page, size int) ([]User, error) {
	var users []User
	db = db.Where("email like ?", emailPattern)
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	offset := size * (page - 1)
	err := db.Offset(offset).Limit(size).Find(&users).Error
	return users, err
}

// // 更新批量用户状态
func UpdateUserStatus(db *gorm.DB, ids []uint, status string) error {
	if err := db.Model(&User{}).Where("id IN ?", ids).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

// // 删除 30 天前未登录
func DeleteInactiveUsers(db *gorm.DB) error {
	cutoff := time.Now().AddDate(0, 0, -30)
	err := db.Where("last_login_at < ?", cutoff).Delete(&User{}).Error

	return err
}

// // scope1：分页
func ScopePage(page, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page < 1 {
			page = 1
		}
		if size < 1 {
			size = 3
		}
		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}

// // scope p2: 18<年龄>30
func ScopeAge(ageBegin, ageEnd int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("age between ? and ?", ageBegin, ageEnd)
	}
}

var users = []User{
	{Name: "Alice", Email: "alice@example.com", Age: 28, Status: "active"},
	{Name: "Alice1", Email: "alice1@example.com", Age: 28, Status: "active"},
	{Name: "Alice2", Email: "alice2@example.com", Age: 28, Status: "inactive"},
	{Name: "Alice3", Email: "alice3@example.com", Age: 28, Status: "inactive"},
	{Name: "Bob", Email: "bob@example.com", Age: 31, Status: "active"},
	{Name: "Bob1", Email: "bob1@example.com", Age: 32, Status: "active"},
	{Name: "Bob2", Email: "bob2@example.com", Age: 33, Status: "inactive"},
	{Name: "Bob3", Email: "bob3example.com", Age: 34, Status: "inactive"},
}
