package participant

import (
	"tupulung/entities"
	"tupulung/entities/web"

	"gorm.io/gorm"
)

type ParticipantRepository struct {
	db *gorm.DB
}

func NewParticipantRepository(db *gorm.DB) ParticipantRepository {
	return ParticipantRepository{
		db: db,
	}
}

func (repo ParticipantRepository) Append(user entities.User, event entities.Event) error {

	err := repo.db.Model(&user).Association("Events").Find(&event)
	if err == nil {
		return web.WebError{Code: 400, Message: "you have joined this event"}
	}
	tx := repo.db.Model(&user).Association("Events").Append(&event)
	if tx != nil {
		return web.WebError{Code: 500, Message: "server error"}
	}
	return nil
}

func (repo ParticipantRepository) Delete(user entities.User, event entities.Event) error {

	err := repo.db.Model(&user).Association("Events").Find(&event)
	if err == nil {
		return web.WebError{Code: 400, Message: "you haven't joined this event"}
	}
	tx := repo.db.Model(&user).Association("Events").Delete(&event)
	if tx != nil {
		return web.WebError{Code: 500, Message: "server error"}
	}
	return nil
}

// 	user, _ := userRepository.Find(6)
// event, _ := eventRepository.Find(2)
// db.Model(&user).Association("Events").Append(&event)
// db.Model(&user).Association("Events").Delete(&event)

// err := db.Model(&user).Association("Events").Find(&event)
// if err != nil {
//     fmt.Println(err.Error())
// } else {
//     fmt.Println("Ada")
// }
