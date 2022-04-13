package comment

import (
	"tupulung/entities"
	"tupulung/entities/web"
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

/*
 * Load pagination data
 * -------------------------------
 * Mengambil data pagination comment berdasarkan filters
 */
func (service CommentService) GetPagination(page, limit int, filters []map[string]string) (web.Pagination, error) {
	totalRows, err := service.commentRepo.CountAll(filters)
	if err != nil {
		return web.Pagination{}, err
	}
	var totalPages int64 = 1
	if limit > 0 {
		totalPages = totalRows / int64(limit)
	}
	if totalPages <= 0 {
		totalPages = 1
	}
	return web.Pagination{
		Page: page,
		Limit: limit,
		TotalPages: int(totalPages),
	}, nil
}