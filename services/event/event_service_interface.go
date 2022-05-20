package event

import (
	"mime/multipart"
	"tupulung/entities"
	web "tupulung/entities/web"
	storageProvider "tupulung/utilities/storage"
)

type EventServiceInterface interface {
	FindAll(limit, page int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.EventResponse, error)
	GetPagination(limit, page int, filters []map[string]string) (web.Pagination, error)
	Find(id int) (entities.EventResponse, error)
	Create(eventRequest entities.EventRequest, userID int, cover *multipart.FileHeader, storageProvider storageProvider.StorageInterface) (entities.EventResponse, error)
	Update(eventRequest entities.EventRequest, id int, userID int, cover *multipart.FileHeader, storageProvider storageProvider.StorageInterface) (entities.EventResponse, error)
	Delete(id int, userID int, storageProvider storageProvider.StorageInterface) error
}
