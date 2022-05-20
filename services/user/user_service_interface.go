package user

import (
	"mime/multipart"
	entity "tupulung/entities"
	storageProvider "tupulung/utilities/storage"
)

type UserServiceInterface interface {
	Find(id int) (entity.UserResponse, error)
	GetJoinedEvents(userID int) ([]entity.EventResponse, error)
	Create(userRequest entity.UserRequest, avatar *multipart.FileHeader, storastorageProvider storageProvider.StorageInterface) (entity.AuthResponse, error)
	Update(userRequest entity.UserRequest, userID int, avatar *multipart.FileHeader, storastorageProvider storageProvider.StorageInterface) (entity.UserResponse, error)
	Delete(userID int, storastorageProvider storageProvider.StorageInterface) error
}
