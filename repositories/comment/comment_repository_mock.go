package comment

import (
	"time"
	"tupulung/entities"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type CommentRepositoryMock struct {
	Mock *mock.Mock
}

func NewCommentRepositoryMock(mock *mock.Mock) *CommentRepositoryMock {
	return &CommentRepositoryMock{
		Mock: mock,
	}
}

var CommentCollection = []entities.Comment{
	{Model: gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		EventID: 1,
		UserID:  1,
		Comment: "some comment",
	},
	{Model: gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		EventID: 2,
		UserID:  2,
		Comment: "lorem ipsum",
	},
}

func (repo CommentRepositoryMock) FindAll(limit int, offset int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.Comment, error) {
	args := repo.Mock.Called(limit, offset, filters, sorts)
	return args.Get(0).([]entities.Comment), args.Error(1)
}

func (repo CommentRepositoryMock) Find(id int) (entities.Comment, error) {
	args := repo.Mock.Called()
	return args.Get(0).(entities.Comment), args.Error(1)
}

func (repo CommentRepositoryMock) FindBy(field string, value string) (entities.Comment, error) {
	args := repo.Mock.Called()
	return args.Get(0).(entities.Comment), args.Error(1)
}

func (repo CommentRepositoryMock) CountAll(filters []map[string]string) (int64, error) {
	args := repo.Mock.Called(filters)
	return int64(args.Int(0)), args.Error(1)
}

func (repo CommentRepositoryMock) Store(comment entities.Comment) (entities.Comment, error) {
	args := repo.Mock.Called()
	return args.Get(0).(entities.Comment), args.Error(1)
}

func (repo CommentRepositoryMock) Update(comment entities.Comment, id int) (entities.Comment, error) {
	args := repo.Mock.Called()
	return args.Get(0).(entities.Comment), args.Error(1)
}

func (repo CommentRepositoryMock) Delete(id int) error {
	args := repo.Mock.Called()
	return args.Error(0)
}
