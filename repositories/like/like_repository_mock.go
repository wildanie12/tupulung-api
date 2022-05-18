package like

import (
	"tupulung/entities"

	"github.com/stretchr/testify/mock"
)

type LikeRepositoryMock struct {
	Mock *mock.Mock
}

func NewLikeRepositoryMock(mock *mock.Mock) *LikeRepositoryMock {
	return &LikeRepositoryMock{
		Mock: mock,
	}
}

var LikeCollection = []entities.Like{
	{
		ID:      1,
		UserID:  1,
		EventID: 1,
	},
	{
		ID:      2,
		UserID:  2,
		EventID: 1,
	},
}

func (repo LikeRepositoryMock) CountLikeByEvent(eventId int) (int64, error) {
	args := repo.Mock.Called()
	return int64(args.Int(0)), args.Error(1)
}

func (repo LikeRepositoryMock) Append(user entities.User, event entities.Event) error {
	args := repo.Mock.Called()
	return args.Error(0)
}

func (repo LikeRepositoryMock) Delete(user entities.User, event entities.Event) error {
	args := repo.Mock.Called()
	return args.Error(0)
}
