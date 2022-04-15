package like

import (
	"reflect"
	"tupulung/entities"
	"tupulung/entities/web"
	eventRepository "tupulung/repositories/event"
	likeRepository "tupulung/repositories/like"
	userRepository "tupulung/repositories/user"

	"github.com/golang-jwt/jwt"
)

type LikeService struct {
	likeRepo  likeRepository.LikeRepositoryInterface
	userRepo  userRepository.UserRepositoryInterface
	eventRepo eventRepository.EventRepositoryInterface
}

func NewLikeService(repository likeRepository.LikeRepository, userRepository userRepository.UserRepository, eventRepository eventRepository.EventRepository) *LikeService {
	return &LikeService{
		likeRepo:  repository,
		userRepo:  userRepository,
		eventRepo: eventRepository,
	}
}

func (service LikeService) Append(token interface{}, eventID int) error {
	user := entities.User{}
	event := entities.Event{}

	tokenID := token.(*jwt.Token)
	claims := tokenID.Claims.(jwt.MapClaims)
	userIDReflect := reflect.ValueOf(claims).MapIndex(reflect.ValueOf("userID"))
	if reflect.ValueOf(userIDReflect.Interface()).Kind().String() != "float64" {
		return web.WebError{Code: 400, Message: "Invalid token, no userdata present"}
	}

	// get user data
	user, err := service.userRepo.Find(int(claims["userID"].(float64)))
	if err != nil {
		return web.WebError{Code: 400, Message: "No user matched with this authenticated user"}
	}

	event, eventErr := service.eventRepo.Find(eventID)
	if eventErr != nil {
		return web.WebError{Code: 400, Message: "Event is not exist"}
	}
	tx := service.likeRepo.Append(user, event)
	if tx != nil {
		return web.WebError{Code: 400, Message: "You are already like this event"}
	}
	return nil
}

func (service LikeService) Delete(token interface{}, eventID int) error {
	user := entities.User{}
	event := entities.Event{}

	tokenID := token.(*jwt.Token)
	claims := tokenID.Claims.(jwt.MapClaims)
	userIDReflect := reflect.ValueOf(claims).MapIndex(reflect.ValueOf("userID"))
	if reflect.ValueOf(userIDReflect.Interface()).Kind().String() != "float64" {
		return web.WebError{Code: 400, Message: "Invalid token, no userdata present"}
	}

	// get user data
	user, err := service.userRepo.Find(int(claims["userID"].(float64)))
	if err != nil {
		return web.WebError{Code: 400, Message: "No user matched with this authenticated user"}
	}

	event, eventErr := service.eventRepo.Find(eventID)
	if eventErr != nil {
		return web.WebError{Code: 400, Message: "Event is not exist"}
	}
	tx := service.likeRepo.Delete(user, event)
	if tx != nil {
		return web.WebError{Code: 400, Message: "You haven't like this event"}
	}
	return nil
}
