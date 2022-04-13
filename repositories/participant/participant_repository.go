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

	participants := []entities.Participant{}
	err := repo.db.Model(entities.Participant{}).Where("user_id = ?", user.ID).Where("event_id = ?", event.ID).Find(&participants)
	if err.RowsAffected > 0 {
		return web.WebError{Code: 400, Message: "you have joined this event"}
	}
	joins := entities.Participant{}
	joins.UserID = user.ID
	joins.EventID = event.ID
	tx := repo.db.Create(&joins)
	if tx.RowsAffected == 0 {
		return web.WebError{Code: 500, Message: "server error"}
	}
	return nil
}

func (repo ParticipantRepository) Delete(user entities.User, event entities.Event) error {

	participants := []entities.Participant{}
	err := repo.db.Model(entities.Participant{}).Where("user_id = ?", user.ID).Where("event_id = ?", event.ID).Find(&participants)
	if err.RowsAffected == 0 {
		return web.WebError{Code: 400, Message: "you haven't joined this event"}
	}
	id := participants[0].ID
	tx := repo.db.Delete(&entities.Participant{}, id)
	if tx.RowsAffected == 0 {
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
