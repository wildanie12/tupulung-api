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
		events, err := eventService.FindAll(0, 0, nil, nil)
		assert.Nil(t, err)
		assert.Equal(t, "Fashion Week", events[0].Title)
	})
}

type mockEventRepositorySuccess struct{}

func (m mockEventRepositorySuccess) FindAll(limit int, offset int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.Event, error)
