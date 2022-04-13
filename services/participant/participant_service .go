package participant

import (
	"reflect"
	"tupulung/entities"
	"tupulung/entities/web"
	eventRepository "tupulung/repositories/event"
	participantRepository "tupulung/repositories/participant"
	userRepository "tupulung/repositories/user"

	"github.com/golang-jwt/jwt"
)

type ParticipantService struct {
	participantRepo participantRepository.ParticipantRepositoryInterface
	userRepo        userRepository.UserRepositoryInterface
	eventRepo       eventRepository.EventRepositoryInterface
}

func NewParticipantService(repository participantRepository.ParticipantRepository, userRepository userRepository.UserRepository, eventRepository eventRepository.EventRepository) *ParticipantService {
	return &ParticipantService{
		participantRepo: repository,
		userRepo:        userRepository,
		eventRepo:       eventRepository,
	}
}

func (service ParticipantService) Append(token interface{}, eventID int) error {
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

	event, _ = service.eventRepo.Find(eventID)
	tx := service.participantRepo.Append(user, event)
	if tx != nil {
		return web.WebError{Code: 400, Message: "You are already join this event"}
	}
	return nil
}

func (service ParticipantService) Delete(token interface{}, eventID int) error {
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

	event, _ = service.eventRepo.Find(eventID)
	tx := service.participantRepo.Delete(user, event)
	if tx != nil {
		return web.WebError{Code: 400, Message: "You are not a member of this event"}
	}
	return nil
}
