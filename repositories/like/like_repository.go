package like

import (
	"tupulung/entities"
	"tupulung/entities/web"

	"gorm.io/gorm"
)

type LikeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) LikeRepository {
	return LikeRepository{
		db: db,
	}
}

func (repo LikeRepository) CountLikeByEvent(eventId int) (int64, error) {
	var count int64
	tx := repo.db.Model(&entities.Like{}).Where("event_id = ?", eventId).Count(&count)
	if tx.Error != nil {
		return -1, web.WebError{Code: 400, Message: tx.Error.Error()}
	}
	return count, nil
}

func (repo LikeRepository) Append(user entities.User, event entities.Event) error {

	likes := []entities.Like{}
	err := repo.db.Model(entities.Like{}).Where("user_id = ?", user.ID).Where("event_id = ?", event.ID).Find(&likes)
	if err.RowsAffected > 0 {
		return web.WebError{Code: 400, Message: "you have liked this event"}
	}
	liker := entities.Like{}
	liker.UserID = user.ID
	liker.EventID = event.ID
	tx := repo.db.Create(&liker)
	if tx.RowsAffected == 0 {
		return web.WebError{Code: 500, Message: "server error"}
	}
	return nil
}

func (repo LikeRepository) Delete(user entities.User, event entities.Event) error {

	likes := []entities.Like{}
	err := repo.db.Model(entities.Like{}).Where("user_id = ?", user.ID).Where("event_id = ?", event.ID).Find(&likes)
	if err.RowsAffected == 0 {
		return web.WebError{Code: 400, Message: "you haven't liked this event"}
	}
	id := likes[0].ID
	tx := repo.db.Delete(&entities.Like{}, id)
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
