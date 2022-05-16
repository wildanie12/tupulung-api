package category

import (
	"tupulung/entities"
	web "tupulung/entities/web"
)

type CategoryServiceInterface interface {
	FindAll(limit, page int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.CategoryResponse, error)
	GetPagination(page, limit int, filters []map[string]string) (web.Pagination, error)
	Find(id int) (entities.CategoryResponse, error)
	Create(categoryRequest entities.CategoryRequest) (entities.CategoryResponse, error)
	Update(categoryRequest entities.CategoryRequest, id int) (entities.CategoryResponse, error)
	Delete(id int) error
}
