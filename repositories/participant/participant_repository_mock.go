package participant

import (
	"tupulung/entities"

	"github.com/stretchr/testify/mock"
)

type ParticipantRepositoryMock struct {
	Mock *mock.Mock
}

func NewParticipantRepositoryMock(mock *mock.Mock) *ParticipantRepositoryMock {
	return &ParticipantRepositoryMock{
		Mock: mock,
	}
}

var ParticipantCollection = []entities.Participant{
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

func (repo ParticipantRepositoryMock) Append(user entities.User, event entities.Event) error {
	args := repo.Mock.Called()
	return args.Error(0)
}

func (repo ParticipantRepositoryMock) Delete(user entities.User, event entities.Event) error {
	args := repo.Mock.Called()
	return args.Error(0)
}
