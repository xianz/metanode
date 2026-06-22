package services

import (
	"blog/models"
	"blog/utils"
	"net/http"

	"gorm.io/gorm"
)

type ListArticlesResponse struct {
	Posts        []models.PostResponse `json:"posts"`
	RowsAffected int                   `json:"rows_affected"`
}

type BlogService struct {
	Db *gorm.DB
}

func NewBlogService(db *gorm.DB) *BlogService {
	return &BlogService{
		Db: db,
	}
}

// 添加
func (bs *BlogService) CreateArticle(req models.PostRequest) (int64, error) {
	result := bs.Db.Create(&models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  req.UserID,
	})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// 更新
func (bs *BlogService) UpdateArticle(id int64, req models.PostRequest) (int64, error) {
	result := bs.Db.Model(&models.Post{}).Where("id = ?", id).Updates(req)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// 删除
func (bs *BlogService) DeleteArticle(id int64) (int64, error) {
	result := bs.Db.Delete(&models.Post{}, id)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// 详细
func (bs *BlogService) GetArticle(id int64) (*models.Post, error) {
	var post models.Post
	if err := bs.Db.Where("id = ?", id).Preload("Comments").First(&post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.NewAppError(http.StatusNotFound, "文章不存在")
		}
		return nil, err
	}
	return &post, nil
}

// 列表
func (bs *BlogService) ListArticles(pageSize, page int) (any, error) {
	var rs []models.PostResponse
	result := bs.scopePage(pageSize, page).Model(&models.Post{}).Scan(&rs)
	if result.Error != nil {
		return nil, result.Error
	}
	// data, _ := json.Marshal(rs)
	// log.Printf("!!!!!%+v\n", string(data))
	return ListArticlesResponse{
		Posts:        rs,
		RowsAffected: int(result.RowsAffected),
	}, nil
}

func (bs *BlogService) scopePage(size, page int) *gorm.DB {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 5
	}
	offset := (page - 1) * size
	return bs.Db.Offset(int(offset)).Limit(size)
}
