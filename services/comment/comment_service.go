package comment

import (
	"tupulung/entities"
	commentRepo "tupulung/repositories/comment"

	"github.com/jinzhu/copier"
)

type CommentService struct {
	commentRepo commentRepo.CommentRepositoryInterface
}

func NewCommentService(commentRepo commentRepo.CommentRepositoryInterface) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
	}
}

/*
 * Find All
 * -------------------------------
 * Mengambil data comment berdasarkan filters dan sorts
 */
func (service CommentService) FindAll(limit, page int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.CommentResponse, error) {

	offset := (page - 1) * limit

	// Repository action find all comment
	comments, err := service.commentRepo.FindAll(limit, offset, filters, sorts)
	if err != nil {
		return []entities.CommentResponse{}, err
	}

	// Konversi ke comment response
	commentsRes := []entities.CommentResponse{}
	copier.Copy(&commentsRes, &comments)
	return commentsRes, nil
}