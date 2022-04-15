package event

import (
	"tupulung/entities"
	"tupulung/entities/web"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return EventRepository{
		db: db,
	}
}

func (repo EventRepository) FindAll(limit int, offset int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.Event, error) {
	events := []entities.Event{}

	builder := repo.db.Preload("User").Preload("Category").Preload("Participants").Preload("Likes").Limit(limit).Offset(offset)
	// Where filters
	for _, filter := range filters {
		builder.Where(filter["field"]+" "+filter["operator"]+" ?", filter["value"])
	}
	// OrderBy Filters
	for _, sort := range sorts {
		builder.Order(clause.OrderByColumn{Column: clause.Column{Name: sort["field"].(string)}, Desc: sort["desc"].(bool)})
	}
	tx := builder.Find(&events)
	if tx.Error != nil {
		return []entities.Event{}, web.WebError{Code: 500, Message: tx.Error.Error()}
	}
	return events, nil
}
func (repo EventRepository) CountAll(filters []map[string]string) (int64, error) {
	var count int64
	builder := repo.db.Model(&entities.Event{})
	// Where filters
	for _, filter := range filters {
		builder.Where(filter["field"]+" "+filter["operator"]+" ?", filter["value"])
	}
	tx := builder.Count(&count)
	if tx.Error != nil {
		return -1, web.WebError{Code: 400, Message: tx.Error.Error()}
	}
	return count, nil
}

func (repo EventRepository) Find(id int) (entities.Event, error) {
	event := entities.Event{}
	tx := repo.db.Preload("User").Preload("Category").Preload("Participants").Preload("Likes").Find(&event, id)
	if tx.Error != nil {
		return entities.Event{}, web.WebError{Code: 500, Message: "server error"}
	} else if tx.RowsAffected <= 0 {
		return entities.Event{}, web.WebError{Code: 400, Message: "cannot get event data with specified id"}
	}
	return event, nil
}

func (repo EventRepository) FindBy(field string, value string) (entities.Event, error) {
	event := entities.Event{}
	tx := repo.db.Preload("User").Preload("Category").Preload("Participants").Preload("Likes").Where(field+" = ?", value).Find(&event)
	if tx.Error != nil {
		return entities.Event{}, web.WebError{Code: 500, Message: tx.Error.Error()}
	} else if tx.RowsAffected <= 0 {
		return entities.Event{}, web.WebError{Code: 400, Message: "The requested ID doesn't match with any record"}
	}
	return event, nil
}

func (repo EventRepository) Store(event entities.Event) (entities.Event, error) {

	tx := repo.db.Preload("User").Preload("Category").Create(&event)
	if tx.Error != nil {
		return entities.Event{}, web.WebError{Code: 500, Message: tx.Error.Error()}
	}
	return event, nil
}

func (repo EventRepository) Update(event entities.Event, id int) (entities.Event, error) {
	tx := repo.db.Save(&event)
	if tx.Error != nil {
		return entities.Event{}, web.WebError{Code: 500, Message: tx.Error.Error()}
	}
	return event, nil
}

func (repo EventRepository) Delete(id int) error {
	tx := repo.db.Delete(&entities.Event{}, id)
	if tx.Error != nil {
		return web.WebError{Code: 500, Message: tx.Error.Error()}
	}
	return nil
}
