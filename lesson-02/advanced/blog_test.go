package advanced

import (
	"testing"

	"lesson02/testutil"

	"gorm.io/gorm"
)

func TestAssociationsBlog(t *testing.T) {
	db := testutil.NewTestDB(t, "blog.db")

	if err := db.AutoMigrate(&User{}, &Post{}, &Comment{}, &Tag{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	db.Exec("PRAGMA foreign_keys = OFF")
	db.Exec("DELETE FROM posts")
	db.Exec("DELETE FROM tags")
	db.Exec("DELETE FROM post_tags")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM comments")
	db.Exec("PRAGMA foreign_keys = ON")

	tags := []Tag{
		{Name: "科技"},
		{Name: "时事"},
	}
	err := db.Create(&tags).Error
	if err != nil {
		t.Fatal(err)
	}

	users := []User{
		{
			Name: "alice",
			Posts: []Post{
				{Title: "标题1", Content: "内容1", Tags: []Tag{tags[1]}, Comments: []Comment{{Content: "评论1"}}},
				{Title: "标题2", Content: "内容2", Tags: []Tag{tags[0]}, Comments: []Comment{{Content: "评论2"}}},
			},
		},
		{
			Name: "bob",
		},
	}
	// 全量保存关联
	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Create(&users).Error; err != nil {
		t.Fatal(err)
	}
	// 修改文章为"标题1" 的作者为bob
	if err := db.Model(&Post{}).Where("title=?", "标题1").Update("UserID", users[1].ID).Error; err != nil {
		t.Fatal(err)
	}

	// 给第一个文章添加一条评论
	comment := Comment{
		Content: "评论3",
	}
	err = db.Model(&users[0].Posts[0]).Association("Comments").Append(&comment)
	if err != nil {
		t.Fatal(err)
	}

	// 查询bob的最新文章
	_, err = GetUserLatestPosts(db, users[1].ID)
	if err != nil {
		t.Fatal(err)
	}
	// t.Logf("posts: %+v", posts)

	// 统计评论数量
	_, err = GetPostsWithCommentCount(db)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("postsWithCommentCount: %+v", postsWithCommentCount)

	// 事务
	post := Post{
		Title:   "标题3",
		Content: "内容3",
		UserID:  users[1].ID,
	}
	if err := PublishPostWithTags(db, &post, []uint{tags[0].ID, tags[1].ID}); err != nil {
		t.Fatal(err)
	}

	//// 软删除评论

	// // 添加软删除功能字段
	// TODO 需要重新设计Comment结构体?
	// type Comment struct {
	// 	ID        uint
	// 	PostID    uint
	// 	Content   string
	// 	CreatedAt time.Time
	// 	DeletedAt gorm.DeletedAt
	// }
	// if err := db.AutoMigrate(&Comment{}); err != nil {
	// 	t.Fatal(err)
	// }
	// 软删除评论
	if err := SoftDeleteComment(db, comment.ID); err != nil {
		t.Fatal(err)
	}
	db.Unscoped().Find(&Comment{}, comment.ID)

	// 彻底删除评论
	if err := HardDeleteComment(db, comment.ID); err != nil {
		t.Fatal(err)
	}
	db.Unscoped().Find(&Comment{}, comment.ID)

}

// // 查询用户最新文章：查询用户发表的最新 10 篇文章，并且包含标签信息
func GetUserLatestPosts(db *gorm.DB, userID uint) ([]Post, error) {
	var posts []Post
	if err := db.Preload("Tags"). // 结构属性名
					Where("user_id=?", userID).      // 条件
					Order("created_at desc").        // 排序
					Limit(10).                       // 限制
					Find(&posts).Error; err != nil { // 查询
		return nil, err
	}
	return posts, nil
}

// // 统计评论数量：使用 Preload + Count 统计每篇文章的评论数量
type PostWithCount struct {
	Title        string
	CommentCount int
}

var postsWithCommentCount []PostWithCount

func GetPostsWithCommentCount(db *gorm.DB) ([]PostWithCount, error) {
	err := db.Model(&Post{}).Select("posts.title as title, COUNT(cm.id) as CommentCount").
		Joins("LEFT JOIN comments cm ON cm.post_id = posts.id").
		Group("posts.title").Scan(&postsWithCommentCount).Error
	if err != nil {
		return nil, err
	}
	return postsWithCommentCount, nil
}

// //在事务中实现文章发布 + 标签绑定
func PublishPostWithTags(db *gorm.DB, post *Post, tagIDs []uint) error {
	// 使用自动事务
	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. 创建文章
		err := tx.Create(post).Error
		if err != nil {
			return err
		}

		// 2. 绑定标签（Many to Many）
		var tags []Tag
		err = tx.Where("id IN ?", tagIDs).Find(&tags).Error
		if err != nil {
			return err
		}
		err = tx.Model(post).Association("Tags").Append(&tags)
		if err != nil {
			return err
		}

		// 3. 更新用户文章数量
		// TODO 方法存疑，操作好像不够简洁
		type User struct {
			PostCount uint
		}
		err = tx.AutoMigrate(&User{})
		if err != nil {
			return err
		}
		type UserPostCount struct {
			UserID    uint
			PostCount uint
		}
		var counts []UserPostCount
		err = tx.Model(&Post{}).Select("user_id, COUNT(*) as post_count").Group("user_id").Scan(&counts).Error
		if err != nil {
			return err
		}
		for _, count := range counts {
			if result := tx.Model(&User{}).Where("id = ?", count.UserID).Update("PostCount", count.PostCount); result.Error != nil {
				return result.Error
			}
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// // 为评论新增软删除，提供"彻底清除"功能
// 软删除评论
func SoftDeleteComment(db *gorm.DB, commentID uint) error {
	err := db.Delete(&Comment{}, commentID).Error
	if err != nil {
		return err
	}
	return nil
}

// 彻底删除评论
func HardDeleteComment(db *gorm.DB, commentID uint) error {
	err := db.Unscoped().Delete(&Comment{}, commentID).Error
	if err != nil {
		return err
	}
	return nil
}
