package participant

import (
	"tupulung/entities"
	"tupulung/entities/web"
	eventRepository "tupulung/repositories/event"
	participantRepository "tupulung/repositories/participant"
	userRepository "tupulung/repositories/user"
)

type ParticipantService struct {
	participantRepo participantRepository.ParticipantRepositoryInterface
	userRepo        userRepository.UserRepositoryInterface
	eventRepo       eventRepository.EventRepositoryInterface
}

func NewParticipantService(repository participantRepository.ParticipantRepositoryInterface,
	userRepository userRepository.UserRepositoryInterface,
	eventRepository eventRepository.EventRepositoryInterface,
) *ParticipantService {
	return &ParticipantService{
		participantRepo: repository,
		userRepo:        userRepository,
		eventRepo:       eventRepository,
	}
}

func (service ParticipantService) Append(userID, eventID int) error {
	user := entities.User{}
	event := entities.Event{}

	// get user data
	user, err := service.userRepo.Find(userID)
	if err != nil {
		return web.WebError{Code: 400, Message: "No user matched with this authenticated user"}
	}

	event, eventErr := service.eventRepo.Find(eventID)
	if eventErr != nil {
		return web.WebError{Code: 400, Message: "Event is not exist"}
	}
	tx := service.participantRepo.Append(user, event)
	if tx != nil {
		return web.WebError{Code: 400, Message: "You are already join this event"}
	}
	return nil
}

func (service ParticipantService) Delete(userID, eventID int) error {
	user := entities.User{}
	event := entities.Event{}

	// get user data
	user, err := service.userRepo.Find(userID)
	if err != nil {
		return web.WebError{Code: 400, Message: "No user matched with this authenticated user"}
	}

	event, eventErr := service.eventRepo.Find(eventID)
	if eventErr != nil {
		return web.WebError{Code: 400, Message: "Event is not exist"}
	}
	tx := service.participantRepo.Delete(user, event)
	if tx != nil {
		return web.WebError{Code: 400, Message: "You are not a member of this event"}
	}
	return nil
}
