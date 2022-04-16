package event

import (
	"mime/multipart"
	"tupulung/entities"
	web "tupulung/entities/web"
)

type EventServiceInterface interface {
	FindAll(limit, page int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.EventResponse, error)
	GetPagination(limit, page int, filters []map[string]string) (web.Pagination, error)
	Find(id int) (entities.EventResponse, error)
	Create(eventRequest entities.EventRequest, tokenReq interface{}, cover *multipart.FileHeader) (entities.EventResponse, error)
	Update(eventRequest entities.EventRequest, id int, tokenReq interface{}, cover *multipart.FileHeader) (entities.EventResponse, error)
	Delete(id int, tokenReq interface{}) error
}
