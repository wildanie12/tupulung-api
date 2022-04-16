package event

import (
	"testing"
	"tupulung/entities"
	"tupulung/repositories/user"

	"github.com/stretchr/testify/assert"
)

func TestFindAll(t *testing.T) {
	t.Run("TestFindAllSuccess", func(t *testing.T) {
		eventService := NewEventService(mockEventRepositorySuccess{}, user.UserRepository{})
		events, err := eventService.FindAll(0, 0, []map[string]string{}, []map[string]interface{}{})
		assert.Nil(t, err)
		assert.Equal(t, "Fashion Week", events[0].Title)
	})
}
func TestFind(t *testing.T) {
	t.Run("TestFindSuccess", func(t *testing.T) {
		eventService := NewEventService(mockEventRepositorySuccess{}, user.UserRepository{})
		events, err := eventService.Find(1)
		assert.Nil(t, err)
		assert.Equal(t, "Fashion Week", events.Title)
	})
}

func TestCreate(t *testing.T) {
	t.Run("TestCreateSuccess", func(t *testing.T) {
		eventService := NewEventService(mockEventRepositorySuccess{}, user.UserRepository{})
		events, err := eventService.Create(entities.EventRequest{Title: "Fashion Week"}, "!", "?")
		assert.Nil(t, err)
		assert.Equal(t, "Fashion Week", events.Title)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("TestUpdateSuccess", func(t *testing.T) {
		eventService := NewEventService(mockEventRepositorySuccess{}, user.UserRepository{})
		events, err := eventService.Update(entities.EventRequest{Title: "Fashion Show"}, 1, "!", "?")
		assert.Nil(t, err)
		assert.Equal(t, "Fashion SHow", events.Title)
	})
}

func TestDelete(t *testing.T) {
	t.Run("TestDeleteSuccess", func(t *testing.T) {
		eventService := NewEventService(mockEventRepositorySuccess{}, user.UserRepository{})
		err := eventService.Delete(1, "!")
		assert.Nil(t, err)
	})
}

type mockEventRepositorySuccess struct{}

func (m mockEventRepositorySuccess) FindAll(limit int, offset int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.Event, error) {
	return []entities.Event{{Title: "Fashion Week"}}, nil
}
func (m mockEventRepositorySuccess) Find(id int) (entities.Event, error) {
	return entities.Event{Title: "Fashion Week"}, nil
}
func (m mockEventRepositorySuccess) FindBy(field string, value string) (entities.Event, error) {
	return entities.Event{Title: "Fashion Week"}, nil
}
func (m mockEventRepositorySuccess) CountAll(filters []map[string]string) (int64, error) {
	return 1, nil
}
func (m mockEventRepositorySuccess) Store(event entities.Event) (entities.Event, error) {
	return entities.Event{Title: "Fashion Week"}, nil
}
func (m mockEventRepositorySuccess) Update(event entities.Event, id int) (entities.Event, error) {
	return entities.Event{Title: "Fashion Show"}, nil
}
func (m mockEventRepositorySuccess) Delete(id int) error {
	return nil
}
