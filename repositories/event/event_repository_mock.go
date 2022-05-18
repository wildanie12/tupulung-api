package event

import (
	"time"
	"tupulung/entities"
	categoryRepository "tupulung/repositories/category"
	userRepository "tupulung/repositories/user"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type EventRepositoryMock struct {
	Mock *mock.Mock
}

func NewEventRepositoryMock(mock *mock.Mock) *EventRepositoryMock {
	return &EventRepositoryMock{
		Mock: mock,
	}
}

var EventCollection = []entities.Event{
	{
		Model:         gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Title:         "Seminar Pendidikan",
		HostedBy:      "nasrul",
		Cover:         "some cover",
		UserID:        1,
		User:          userRepository.UserCollection[0],
		CategoryID:    1,
		Category:      categoryRepository.CategoryCollection[0],
		DatetimeEvent: time.Now(),
		Location:      "surabaya",
		Description:   "some description",
	},
	{
		Model:         gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Title:         "Seminar Teknologi",
		HostedBy:      "ramadhan",
		Cover:         "some cover",
		UserID:        2,
		User:          userRepository.UserCollection[1],
		CategoryID:    2,
		Category:      categoryRepository.CategoryCollection[1],
		DatetimeEvent: time.Now(),
		Location:      "jakareta",
		Description:   "some description",
	},
}

func (repo EventRepositoryMock) FindAll(limit int, offset int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.Event, error) {
	args := repo.Mock.Called(limit, offset, filters, sorts)
	return args.Get(0).([]entities.Event), args.Error(1)
}
func (repo EventRepositoryMock) Find(id int) (entities.Event, error) {
	args := repo.Mock.Called()
	return args.Get(0).(entities.Event), args.Error(1)
}
func (repo EventRepositoryMock) FindBy(field string, value string) (entities.Event, error) {
	args := repo.Mock.Called()
	return args.Get(0).(entities.Event), args.Error(1)
}
func (repo EventRepositoryMock) CountAll(filters []map[string]string) (int64, error) {
	args := repo.Mock.Called(filters)
	return int64(args.Int(0)), args.Error(1)
}
func (repo EventRepositoryMock) Store(user entities.Event) (entities.Event, error) {
	args := repo.Mock.Called()
	return args.Get(0).(entities.Event), args.Error(1)
}
func (repo EventRepositoryMock) Update(user entities.Event, id int) (entities.Event, error) {
	args := repo.Mock.Called()
	return args.Get(0).(entities.Event), args.Error(1)
}
func (repo EventRepositoryMock) Delete(id int) error {
	args := repo.Mock.Called()
	return args.Error(0)
}
func (repo EventRepositoryMock) DeleteBatch(filters []map[string]string) error {
	args := repo.Mock.Called()
	return args.Error(0)
}
