package event

import (
	"mime/multipart"
	"net/url"
	"reflect"
	"strings"
	"time"
	"tupulung/deliveries/helpers"
	"tupulung/entities"

	web "tupulung/entities/web"
	eventRepository "tupulung/repositories/event"
	userRepository "tupulung/repositories/user"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type EventService struct {
	eventRepo eventRepository.EventRepositoryInterface
	userRepo  userRepository.UserRepositoryInterface
}

func NewEventService(repository eventRepository.EventRepositoryInterface, userRepository userRepository.UserRepositoryInterface) *EventService {
	return &EventService{
		eventRepo: repository,
		userRepo:  userRepository,
	}
}

/*
 * --------------------------
 * Get List of event
 * --------------------------
 */
func (service EventService) FindAll(limit, page int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.EventResponse, error) {

	offset := (page - 1) * limit

	eventsRes := []entities.EventResponse{}
	events, err := service.eventRepo.FindAll(limit, offset, filters, sorts)
	copier.Copy(&eventsRes, &events)
	return eventsRes, err
}

/*
 * --------------------------
 * Load pagination data
 * --------------------------
 */
func (service EventService) GetPagination(limit, page int, filters []map[string]string) (web.Pagination, error) {
	totalRows, err := service.eventRepo.CountAll(filters)
	if err != nil {
		return web.Pagination{}, err
	}
	if limit <= 0 {
		limit = 1
	}
	totalPages := totalRows / int64(limit)
	if totalRows%int64(limit) > 0 {
		totalPages++
	}

	return web.Pagination{
		Page:       page,
		Limit:      limit,
		TotalPages: int(totalPages),
	}, nil
}

/*
 * --------------------------
 * Get single event data based on ID
 * --------------------------
 */
func (service EventService) Find(id int) (entities.EventResponse, error) {

	event, err := service.eventRepo.Find(id)
	eventRes := entities.EventResponse{}
	copier.Copy(&eventRes, &event)

	return eventRes, err
}

/*
 * --------------------------
 * Create event resource
 * --------------------------
 */
func (service EventService) Create(eventRequest entities.EventRequest, tokenReq interface{}, cover *multipart.FileHeader) (entities.EventResponse, error) {
	// convert event to entities entities
	event := entities.Event{}
	copier.Copy(&event, &eventRequest)

	// Translate token
	token := tokenReq.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userIDReflect := reflect.ValueOf(claims).MapIndex(reflect.ValueOf("userID"))
	if reflect.ValueOf(userIDReflect.Interface()).Kind().String() != "float64" {
		return entities.EventResponse{}, web.WebError{Code: 400, Message: "Invalid token, no userdata present"}
	}

	// get user data
	user, err := service.userRepo.Find(int(claims["userID"].(float64)))
	if err != nil {
		return entities.EventResponse{}, web.WebError{Code: 400, Message: "No user matched with this authenticated user"}
	}
	event.UserID = user.ID

	// repository action
	if eventRequest.DatetimeEvent != "" {
		datetime, err := time.Parse("2006-01-02", eventRequest.DatetimeEvent)
		if err != nil {
			return entities.EventResponse{}, web.WebError{Code: 400, Message: "date time event format is invalid"}
		}
		event.DatetimeEvent = datetime
	}

	if cover != nil {
		coverFile, err := cover.Open()
		if err != nil {
			return entities.EventResponse{}, web.WebError{Code: 500, Message: "Cannot process cover image"}
		}
		defer coverFile.Close()

		// Upload cover to S3
		filename := uuid.New().String() + cover.Filename
		coverURL, err := helpers.UploadFileToS3("event/cover/"+filename, coverFile)
		if err != nil {
			return entities.EventResponse{}, web.WebError{Code: 500, Message: err.Error()}
		}
		event.Cover = coverURL
	}

	event, err = service.eventRepo.Store(event)
	if err != nil {
		return entities.EventResponse{}, err
	}

	// get event data
	eventRes, err := service.Find(int(event.ID))
	if err != nil {
		return entities.EventResponse{}, web.WebError{Code: 500, Message: "Cannot get newly created event"}
	}

	return eventRes, nil
}

/*
 * --------------------------
 * Update event resource
 * --------------------------
 */
func (service EventService) Update(eventRequest entities.EventRequest, id int, tokenReq interface{}, cover *multipart.FileHeader) (entities.EventResponse, error) {

	// Find event
	event, err := service.eventRepo.Find(id)
	if err != nil {
		return entities.EventResponse{}, web.WebError{Code: 400, Message: "The requested ID doesn't match with any record"}
	}

	// Translate token
	token := tokenReq.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userIDReflect := reflect.ValueOf(claims).MapIndex(reflect.ValueOf("userID"))
	if reflect.ValueOf(userIDReflect.Interface()).Kind().String() != "float64" {
		return entities.EventResponse{}, web.WebError{Code: 400, Message: "Invalid token, no userdata present"}
	}

	// get user data
	user, err := service.userRepo.Find(int(claims["userID"].(float64)))
	if err != nil {
		return entities.EventResponse{}, web.WebError{Code: 400, Message: "No user matched with this authenticated user"}
	}
	if event.UserID != user.ID {
		return entities.EventResponse{}, web.WebError{Code: 401, Message: "Cannot update event that belongs to someone else"}
	}
	if eventRequest.DatetimeEvent != "" {
		datetime, err := time.Parse("2006-01-02", eventRequest.DatetimeEvent)
		if err != nil {
			return entities.EventResponse{}, web.WebError{Code: 400, Message: "date time event format is invalid"}
		}
		event.DatetimeEvent = datetime
	}
	if cover != nil {

		// Delete previous cover
		if event.Cover != "" {
			u, _ := url.Parse(event.Cover)
			objectPathS3 := strings.TrimPrefix(u.Path, "/")
			helpers.DeleteFromS3(objectPathS3)
		}

		coverFile, err := cover.Open()
		if err != nil {
			return entities.EventResponse{}, web.WebError{Code: 500, Message: "cannot read cover image file"}
		}
		// Upload cover to S3
		filename := uuid.New().String() + cover.Filename
		coverURL, err := helpers.UploadFileToS3("event/cover/"+filename, coverFile)
		if err != nil {
			return entities.EventResponse{}, web.WebError{Code: 500, Message: err.Error()}
		}
		event.Cover = coverURL
	}
	// Copy request to found event
	copier.CopyWithOption(&event, &eventRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	// repository action
	event, err = service.eventRepo.Update(event, id)
	if err != nil {
		return entities.EventResponse{}, err
	}

	// get event data
	eventRes, err := service.Find(int(event.ID))
	if err != nil {
		return entities.EventResponse{}, web.WebError{Code: 500, Message: "Cannot get newly created event"}
	}

	return eventRes, err
}

/*
 * --------------------------
 * Delete resource data
 * --------------------------
 */
func (service EventService) Delete(id int, tokenReq interface{}) error {
	// Find event
	event, err := service.eventRepo.Find(id)
	if err != nil {
		return web.WebError{Code: 400, Message: "The requested ID doesn't match with any record"}
	}

	// Translate token
	token := tokenReq.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userIDReflect := reflect.ValueOf(claims).MapIndex(reflect.ValueOf("userID"))
	if reflect.ValueOf(userIDReflect.Interface()).Kind().String() != "float64" {
		return web.WebError{Code: 400, Message: "Invalid token, no userdata present"}
	}

	// get user data
	user, err := service.userRepo.Find(int(claims["userID"].(float64)))
	if err != nil {
		return web.WebError{Code: 400, Message: "No user matched with this authenticated user"}
	}
	if event.UserID != user.ID {
		return web.WebError{Code: 401, Message: "Cannot update event that belongs to someone else"}
	}

	// Delete previous cover
	if event.Cover != "" {
		u, _ := url.Parse(event.Cover)
		objectPathS3 := strings.TrimPrefix(u.Path, "/")
		helpers.DeleteFromS3(objectPathS3)
	}

	// Repository action
	err = service.eventRepo.Delete(id)
	return err
}
