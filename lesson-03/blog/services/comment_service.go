package services

import (
	"blog/models"

	"gorm.io/gorm"
)

type CommentService struct {
	Db *gorm.DB
}

func (cs *CommentService) CreateComment(comment map[string]any) (int64, error) {
	result := cs.Db.Model(&models.Comment{}).Create(&comment)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (cs *CommentService) GetCommentsByPostID(postID int64) (any, error) {
	var comments []models.Comment
	result := cs.Db.Where("post_id = ?", postID).Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return models.ListCommentResponse{
		Comments:     comments,
		RowsAffected: int(result.RowsAffected),
	}, nil
}
