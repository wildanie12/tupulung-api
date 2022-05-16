package category

import (
	"time"
	"tupulung/entities"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type CategoryRepositoryMock struct {
	Mock *mock.Mock
}

func NewCategoryRepositoryMock(mock *mock.Mock) *CategoryRepositoryMock {
	return &CategoryRepositoryMock{
		Mock: mock,
	}
}

var CategoryCollection = []entities.Category{
	{Model: gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Title: "Education",
	},
	{Model: gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Title: "Technology",
	},
}

func (repo CategoryRepositoryMock) FindAll(limit int, offset int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.Category, error) {
	args := repo.Mock.Called(limit, offset, filters, sorts)
	return args.Get(0).([]entities.Category), args.Error(1)
}

func (repo CategoryRepositoryMock) Find(id int) (entities.Category, error) {
	args := repo.Mock.Called(id)
	return args.Get(0).(entities.Category), args.Error(1)
}

func (repo CategoryRepositoryMock) CountAll(filters []map[string]string) (int64, error) {
	args := repo.Mock.Called(filters)
	return int64(args.Int(0)), args.Error(1)
}

func (repo CategoryRepositoryMock) Store(category entities.Category) (entities.Category, error) {
	args := repo.Mock.Called()
	return args.Get(0).(entities.Category), args.Error(1)
}

func (repo CategoryRepositoryMock) Update(category entities.Category, id int) (entities.Category, error) {
	args := repo.Mock.Called()
	return args.Get(0).(entities.Category), args.Error(1)
}

func (repo CategoryRepositoryMock) Delete(id int) error {
	args := repo.Mock.Called(id)
	return args.Error(0)
}
