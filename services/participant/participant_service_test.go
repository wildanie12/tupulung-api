package participant_test

import (
	"testing"
	"tupulung/entities"
	"tupulung/entities/web"
	eventRepository "tupulung/repositories/event"
	participantRepository "tupulung/repositories/participant"
	userRepository "tupulung/repositories/user"
	participantService "tupulung/services/participant"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAppend(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userSample := userRepository.UserCollection[0]
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(userSample, nil)
		eventSample := eventRepository.EventCollection[0]
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Find").Return(eventSample, nil)
		participantRepositoryMock := participantRepository.NewParticipantRepositoryMock(&mock.Mock{})
		participantRepositoryMock.Mock.On("Append").Return(nil)

		Service := participantService.NewParticipantService(
			participantRepositoryMock,
			userRepositoryMock,
			eventRepositoryMock,
		)
		err := Service.Append(int(userSample.ID), int(eventSample.ID))
		assert.Nil(t, err)
	})
	t.Run("repo-fail", func(t *testing.T) {
		userSample := userRepository.UserCollection[0]
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(userSample, nil)
		eventSample := eventRepository.EventCollection[0]
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Find").Return(eventSample, nil)
		participantRepositoryMock := participantRepository.NewParticipantRepositoryMock(&mock.Mock{})
		participantRepositoryMock.Mock.On("Append").Return(web.WebError{})
		Service := participantService.NewParticipantService(
			participantRepositoryMock,
			userRepositoryMock,
			eventRepositoryMock,
		)
		err := Service.Append(int(userSample.ID), int(eventSample.ID))
		assert.Error(t, err)
	})
	t.Run("repo-fail-user", func(t *testing.T) {
		userSample := userRepository.UserCollection[0]
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(entities.User{}, web.WebError{})
		eventSample := eventRepository.EventCollection[0]
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Find").Return(eventSample, nil)
		participantRepositoryMock := participantRepository.NewParticipantRepositoryMock(&mock.Mock{})
		participantRepositoryMock.Mock.On("Append").Return(web.WebError{})
		Service := participantService.NewParticipantService(
			participantRepositoryMock,
			userRepositoryMock,
			eventRepositoryMock,
		)
		err := Service.Append(int(userSample.ID), int(eventSample.ID))
		assert.Error(t, err)
	})
	t.Run("repo-fail-event", func(t *testing.T) {
		userSample := userRepository.UserCollection[0]
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(userSample, nil)
		eventSample := eventRepository.EventCollection[0]
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Find").Return(entities.Event{}, web.WebError{})
		participantRepositoryMock := participantRepository.NewParticipantRepositoryMock(&mock.Mock{})
		participantRepositoryMock.Mock.On("Append").Return(web.WebError{})
		Service := participantService.NewParticipantService(
			participantRepositoryMock,
			userRepositoryMock,
			eventRepositoryMock,
		)
		err := Service.Append(int(userSample.ID), int(eventSample.ID))
		assert.Error(t, err)
	})
}
func TestDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userSample := userRepository.UserCollection[0]
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(userSample, nil)
		eventSample := eventRepository.EventCollection[0]
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Find").Return(eventSample, nil)
		participantRepositoryMock := participantRepository.NewParticipantRepositoryMock(&mock.Mock{})
		participantRepositoryMock.Mock.On("Delete").Return(nil)

		Service := participantService.NewParticipantService(
			participantRepositoryMock,
			userRepositoryMock,
			eventRepositoryMock,
		)
		err := Service.Delete(int(userSample.ID), int(eventSample.ID))
		assert.Nil(t, err)
	})
	t.Run("repo-fail", func(t *testing.T) {
		userSample := userRepository.UserCollection[0]
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(userSample, nil)
		eventSample := eventRepository.EventCollection[0]
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Find").Return(eventSample, nil)
		participantRepositoryMock := participantRepository.NewParticipantRepositoryMock(&mock.Mock{})
		participantRepositoryMock.Mock.On("Delete").Return(web.WebError{})
		Service := participantService.NewParticipantService(
			participantRepositoryMock,
			userRepositoryMock,
			eventRepositoryMock,
		)
		err := Service.Delete(int(userSample.ID), int(eventSample.ID))
		assert.Error(t, err)
	})
	t.Run("repo-fail-user", func(t *testing.T) {
		userSample := userRepository.UserCollection[0]
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(entities.User{}, web.WebError{})
		eventSample := eventRepository.EventCollection[0]
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Find").Return(eventSample, nil)
		participantRepositoryMock := participantRepository.NewParticipantRepositoryMock(&mock.Mock{})
		participantRepositoryMock.Mock.On("Delete").Return(web.WebError{})
		Service := participantService.NewParticipantService(
			participantRepositoryMock,
			userRepositoryMock,
			eventRepositoryMock,
		)
		err := Service.Delete(int(userSample.ID), int(eventSample.ID))
		assert.Error(t, err)
	})
	t.Run("repo-fail-event", func(t *testing.T) {
		userSample := userRepository.UserCollection[0]
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(userSample, nil)
		eventSample := eventRepository.EventCollection[0]
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Find").Return(entities.Event{}, web.WebError{})
		participantRepositoryMock := participantRepository.NewParticipantRepositoryMock(&mock.Mock{})
		participantRepositoryMock.Mock.On("Delete").Return(web.WebError{})
		Service := participantService.NewParticipantService(
			participantRepositoryMock,
			userRepositoryMock,
			eventRepositoryMock,
		)
		err := Service.Delete(int(userSample.ID), int(eventSample.ID))
		assert.Error(t, err)
	})
}
