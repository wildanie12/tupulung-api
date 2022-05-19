package comment

import (
	"tupulung/entities"
	"tupulung/entities/web"
)

type CommentServiceInterface interface {
	FindAll(limit, page int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.CommentResponse, error)
	GetPagination(page, limit int, filters []map[string]string) (web.Pagination, error)
	Create(commentRequest entities.CommentRequest, eventID int, tokenReq interface{}) (entities.CommentResponse, error)
	Update(commentRequest entities.CommentRequest, id int, tokenReq interface{}) (entities.CommentResponse, error)
	Delete(id int, tokenReq interface{}) error
}
