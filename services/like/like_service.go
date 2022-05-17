package like

import (
	"tupulung/entities"
	"tupulung/entities/web"
	eventRepository "tupulung/repositories/event"
	likeRepository "tupulung/repositories/like"
	userRepository "tupulung/repositories/user"
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

func (service LikeService) Append(ID int, eventID int) error {
	user := entities.User{}
	event := entities.Event{}

	// get user data
	user, err := service.userRepo.Find(ID)
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

func (service LikeService) Delete(ID int, eventID int) error {
	user := entities.User{}
	event := entities.Event{}

	// get user data
	user, err := service.userRepo.Find(ID)
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
