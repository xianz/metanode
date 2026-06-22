package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model `json:"-"`
	Username   string `gorm:"uniqueIndex"`
	Password   string `gorm:"not null" json:"-"`
	Email      string `gorm:"unique;not null"`
}

type Post struct {
	gorm.Model `json:"-"`
	Title      string `gorm:"not null"`
	Content    string `gorm:"not null"`
	UserID     uint   `gorm:"not null"`
	Comments   []Comment
}

type Comment struct {
	gorm.Model `json:"-"`
	PostID     uint   `gorm:"not null"`
	UserID     uint   `gorm:"not null"`
	Content    string `gorm:"not null"`
}

type PostResponse struct {
	// Post  // 暂时无法通过结构体中加入json注释来隐藏字段序列化
	// // 显式遮蔽 gorm.Model 的四个字段，使其在 JSON 序列化时被完全忽略
	// // ID        uint           //`json:"-"`
	// CreatedAt time.Time      `json:"-"`
	// UpdatedAt time.Time      `json:"-"`
	// DeletedAt gorm.DeletedAt `json:"-"`
	ID      uint
	Title   string
	Content string
	UserID  uint
	// Comments []Comment `json:"comments,omitempty"` // Scan方法用不了
}

type PostRequest struct {
	// ID      uint
	Title   string // `json:"title"`
	Content string //`json:"content"`
	UserID  uint
}

type UserRequest struct {
	Username string `validate:"required,min=3,max=20,alphanum" json:"username" label:"用户名"`
	Password string `validate:"required,min=8,max=32,containsany=!@#$%^&*" json:"password" label:"密码"`
	Email    string `validate:"required,email" json:"email" label:"邮箱"`
}

type ListCommentResponse struct {
	Comments     []Comment `json:"comments"`
	RowsAffected int       `json:"rows_affected"`
}
