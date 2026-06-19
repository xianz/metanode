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

func (bs *BlogService) CreateArticle(req models.ArticleRequest) (*models.Post, error) {
	var post models.Post
	if err := bs.Db.Table("blog_posts").Create(&req).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (bs *BlogService) UpdateArticle() {

}

func (bs *BlogService) DeleteArticle(id int) error {
	return bs.Db.Delete(&models.Post{}, id).Error
}

func (bs *BlogService) GetArticle(id int) (*models.Post, error) {
	var post models.Post
	if err := bs.Db.Where("id = ?", id).Preload("Comments").First(&post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.NewAppError(http.StatusNotFound, "文章不存在")
		}
		return nil, err
	}
	return &post, nil
}

func (bs *BlogService) ListArticles(pageSize, page int) (ListArticlesResponse, error) {
	var rs []models.PostResponse
	result := bs.scopePage(pageSize, page).Table("blog_posts").Scan(&rs)
	if result.Error != nil {
		return ListArticlesResponse{}, result.Error
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
