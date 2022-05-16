package user

import (
	"time"
	"tupulung/entities"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type UserRepositoryMock struct {
	Mock *mock.Mock
}

func NewUserRepositoryMock(mock *mock.Mock) *UserRepositoryMock {
	return &UserRepositoryMock{
		Mock: mock,
	}
}

var UserCollection = []entities.User{
	{
		Model:    gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Name:     "test1",
		Email:    "test1@mail.com",
		Password: "test",
		Gender:   "male",
		DOB:      time.Now(),
		Address:  "jl. reformasi",
		Avatar:   "some avatar",
	},
	{
		Model:    gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Name:     "test2",
		Email:    "test2@mail.com",
		Password: "test",
		Gender:   "male",
		DOB:      time.Now(),
		Address:  "jl. reformasi",
		Avatar:   "some avatar",
	},
}

func (repo UserRepositoryMock) Find(id int) (entities.User, error) {
	args := repo.Mock.Called(id)
	return args.Get(0).(entities.User), args.Error(1)
}
func (repo UserRepositoryMock) GetJoinedEvents(id int) ([]entities.Event, error) {
	args := repo.Mock.Called()
	return args.Get(0).([]entities.Event), args.Error(1)
}
func (repo UserRepositoryMock) FindBy(field string, value string) (entities.User, error) {
	args := repo.Mock.Called()
	return args.Get(0).(entities.User), args.Error(1)
}
func (repo UserRepositoryMock) Store(user entities.User) (entities.User, error) {
	args := repo.Mock.Called()
	return args.Get(0).(entities.User), args.Error(1)
}
func (repo UserRepositoryMock) Update(user entities.User, id int) (entities.User, error) {
	args := repo.Mock.Called()
	return args.Get(0).(entities.User), args.Error(1)
}
func (repo UserRepositoryMock) Delete(id int) error {
	args := repo.Mock.Called()
	return args.Error(0)
}
