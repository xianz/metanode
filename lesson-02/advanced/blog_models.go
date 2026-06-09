package advanced

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint
	Name      string `gorm:"uniqueIndex"`
	Posts     []Post
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Post struct {
	ID        uint
	Title     string
	Content   string
	UserID    uint
	Tags      []Tag `gorm:"many2many:post_tags;"`
	Comments  []Comment
	CreatedAt time.Time
}

type Comment struct {
	ID        uint
	PostID    uint
	Content   string
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type Tag struct {
	ID        uint
	Posts     []Post `gorm:"many2many:post_tags;"`
	Name      string `gorm:"uniqueIndex"`
	CreatedAt time.Time
}
