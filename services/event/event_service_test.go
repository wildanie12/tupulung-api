package event_test

import (
	"mime/multipart"
	"net/textproto"
	"testing"
	"tupulung/entities"
	web "tupulung/entities/web"
	eventRepository "tupulung/repositories/event"
	likeRepository "tupulung/repositories/like"
	userRepository "tupulung/repositories/user"
	eventService "tupulung/services/event"
	_storageProvider "tupulung/utilities/storage"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFindAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		likeRepositoryMock := likeRepository.NewLikeRepositoryMock(&mock.Mock{})
		likeRepositoryMock.Mock.On("CountLikeByEvent").Return(0, nil)
		eventSample := eventRepository.EventCollection
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On(
			"FindAll",
			0, 0,
			[]map[string]string{},
			[]map[string]interface{}{},
		).Return(eventSample, nil)

		service := eventService.NewEventService(
			eventRepositoryMock,
			userRepository.NewUserRepositoryMock(&mock.Mock{}),
			likeRepositoryMock,
		)
		data, err := service.FindAll(0, 0, []map[string]string{}, []map[string]interface{}{})

		// Konversi expected data ke response
		eventRes := []entities.EventResponse{}
		copier.Copy(&eventRes, &eventSample)

		assert.Nil(t, err)
		assert.Equal(t, eventRes, data)
	})
	t.Run("repo-fail", func(t *testing.T) {
		likeRepositoryMock := likeRepository.NewLikeRepositoryMock(&mock.Mock{})
		likeRepositoryMock.Mock.On("CountLikeByEvent").Return(0, nil)
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On(
			"FindAll",
			0, 0,
			[]map[string]string{},
			[]map[string]interface{}{},
		).Return([]entities.Event{}, web.WebError{})

		service := eventService.NewEventService(
			eventRepositoryMock,
			userRepository.NewUserRepositoryMock(&mock.Mock{}),
			likeRepositoryMock,
		)
		data, err := service.FindAll(0, 0, []map[string]string{}, []map[string]interface{}{})

		// Konversi expected data ke response
		eventRes := []entities.EventResponse{}

		assert.Error(t, err)
		assert.Equal(t, eventRes, data)
	})
}

func TestGetPagination(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("CountAll", []map[string]string{}).Return(20, nil)

		Service := eventService.NewEventService(
			eventRepositoryMock,
			userRepository.NewUserRepositoryMock(&mock.Mock{}),
			likeRepository.NewLikeRepositoryMock(&mock.Mock{}),
		)
		actual, err := Service.GetPagination(5, 1, []map[string]string{})

		expected := web.Pagination{
			Page:       1,
			Limit:      5,
			TotalPages: int(4),
		}
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("repo-fail", func(t *testing.T) {
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("CountAll", []map[string]string{}).Return(0, web.WebError{})

		Service := eventService.NewEventService(
			eventRepositoryMock,
			userRepository.NewUserRepositoryMock(&mock.Mock{}),
			likeRepository.NewLikeRepositoryMock(&mock.Mock{}),
		)
		actual, err := Service.GetPagination(5, 1, []map[string]string{})
		assert.Error(t, err)
		assert.Equal(t, web.Pagination{}, actual)
	})
	t.Run("limit-zero", func(t *testing.T) {
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("CountAll", []map[string]string{}).Return(20, nil)

		Service := eventService.NewEventService(
			eventRepositoryMock,
			userRepository.NewUserRepositoryMock(&mock.Mock{}),
			likeRepository.NewLikeRepositoryMock(&mock.Mock{}),
		)
		actual, err := Service.GetPagination(1, 1, []map[string]string{})

		expected := web.Pagination{
			Page:       1,
			Limit:      1,
			TotalPages: int(20),
		}
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("page-zero", func(t *testing.T) {
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("CountAll", []map[string]string{}).Return(1, nil)

		Service := eventService.NewEventService(
			eventRepositoryMock,
			userRepository.NewUserRepositoryMock(&mock.Mock{}),
			likeRepository.NewLikeRepositoryMock(&mock.Mock{}),
		)
		actual, err := Service.GetPagination(1, 0, []map[string]string{})

		expected := web.Pagination{
			Page:       0,
			Limit:      1,
			TotalPages: int(1),
		}
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("added-page-on-active-module", func(t *testing.T) {
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("CountAll", []map[string]string{}).Return(20, nil)

		Service := eventService.NewEventService(
			eventRepositoryMock,
			userRepository.NewUserRepositoryMock(&mock.Mock{}),
			likeRepository.NewLikeRepositoryMock(&mock.Mock{}),
		)
		actual, err := Service.GetPagination(1, 5, []map[string]string{})

		expected := web.Pagination{
			Page:       5,
			Limit:      1,
			TotalPages: int(20),
		}
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestFind(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		likeRepositoryMock := likeRepository.NewLikeRepositoryMock(&mock.Mock{})
		likeRepositoryMock.Mock.On("CountLikeByEvent").Return(0, nil)
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		eventSample := eventRepository.EventCollection[0]
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Find").Return(eventSample, nil)

		Service := eventService.NewEventService(
			eventRepositoryMock,
			userRepositoryMock,
			likeRepositoryMock,
		)
		_, err := Service.Find(int(eventSample.ID))

		assert.Nil(t, err)
	})
	t.Run("failed", func(t *testing.T) {
		likeRepositoryMock := likeRepository.NewLikeRepositoryMock(&mock.Mock{})
		likeRepositoryMock.Mock.On("CountLikeByEvent").Return(0, nil)
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		eventSample := eventRepository.EventCollection[0]
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Find").Return(entities.Event{}, web.WebError{})

		Service := eventService.NewEventService(
			eventRepositoryMock,
			userRepositoryMock,
			likeRepositoryMock,
		)
		_, err := Service.Find(int(eventSample.ID))

		assert.Error(t, err)
	})
}

func TestCreate(t *testing.T) {
	sampleCentral := eventRepository.EventCollection[0]
	sampleUser := userRepository.UserCollection[0]
	sampleRequestCentral := entities.EventRequest{}
	copier.Copy(&sampleRequestCentral, &sampleCentral)
	sampleRequestCentral.DatetimeEvent = "1999-12-12"
	cover := &multipart.FileHeader{
		Filename: "cover.jpg",
		Header: textproto.MIMEHeader{
			"Content-Disposition": []string{
				"form-data; name=\"cover\"; filename=\"cover.png\"",
			},
			"Content-Type": []string{
				"image/png",
			},
		},
		Size: 155 * 1024,
	}
	t.Run("success", func(t *testing.T) {
		sampleEvent := sampleCentral
		sampleRequest := sampleRequestCentral
		sampleFileRequest := cover

		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(sampleUser, nil)
		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("UploadFromRequest").Return("example.com/images.png", nil)
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Store").Return(sampleEvent, nil)
		eventRepositoryMock.Mock.On("Find").Return(sampleEvent, nil)
		likeRepositoryMock := likeRepository.NewLikeRepositoryMock(&mock.Mock{})
		likeRepositoryMock.Mock.On("CountLikeByEvent").Return(0, nil)

		Service := eventService.NewEventService(
			eventRepositoryMock,
			userRepositoryMock,
			likeRepositoryMock,
		)
		actual, err := Service.Create(sampleRequest, int(sampleUser.ID), sampleFileRequest, storageProvider)

		expected := entities.EventResponse{}
		copier.Copy(&expected, &sampleEvent)

		assert.Nil(t, err)
		assert.Equal(t, expected.ID, actual.ID)
	})
	t.Run("validation-fail", func(t *testing.T) {
		sampleEvent := sampleCentral
		sampleRequest := sampleRequestCentral
		sampleFileRequest := cover

		sampleRequest.Title = ""
		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(sampleUser, nil)
		storageProvider.Mock.On("UploadFromRequest").Return("example.com/images.png", nil)
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Store").Return(sampleEvent, nil)
		eventRepositoryMock.Mock.On("Find").Return(sampleEvent, nil)
		likeRepositoryMock := likeRepository.NewLikeRepositoryMock(&mock.Mock{})
		likeRepositoryMock.Mock.On("CountLikeByEvent").Return(0, nil)

		Service := eventService.NewEventService(
			eventRepositoryMock,
			userRepositoryMock,
			likeRepositoryMock,
		)
		actual, err := Service.Create(sampleRequest, int(sampleUser.ID), sampleFileRequest, storageProvider)

		expected := entities.EventResponse{}
		copier.Copy(&expected, &sampleEvent)

		assert.Error(t, err)
		assert.Equal(t, entities.EventResponse{}, actual)
	})
	t.Run("invalid-datetime", func(t *testing.T) {
		sampleEvent := sampleCentral
		sampleRequest := sampleRequestCentral
		sampleFileRequest := cover

		sampleRequest.DatetimeEvent = "2022222222"
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(sampleUser, nil)
		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("UploadFromRequest").Return("example.com/images.png", nil)
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Store").Return(sampleEvent, nil)
		eventRepositoryMock.Mock.On("Find").Return(sampleEvent, nil)
		likeRepositoryMock := likeRepository.NewLikeRepositoryMock(&mock.Mock{})
		likeRepositoryMock.Mock.On("CountLikeByEvent").Return(0, nil)

		Service := eventService.NewEventService(
			eventRepositoryMock,
			userRepositoryMock,
			likeRepositoryMock,
		)
		actual, err := Service.Create(sampleRequest, int(sampleUser.ID), sampleFileRequest, storageProvider)

		expected := entities.EventResponse{}
		copier.Copy(&expected, &sampleEvent)

		assert.Error(t, err)
		assert.Equal(t, entities.EventResponse{}, actual)
	})
	t.Run("upload-fail", func(t *testing.T) {
		sampleEvent := sampleCentral
		sampleRequest := sampleRequestCentral
		sampleFileRequest := cover
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(sampleUser, nil)
		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("UploadFromRequest").Return("", web.WebError{})
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Store").Return(sampleEvent, nil)
		eventRepositoryMock.Mock.On("Find").Return(sampleEvent, nil)
		likeRepositoryMock := likeRepository.NewLikeRepositoryMock(&mock.Mock{})
		likeRepositoryMock.Mock.On("CountLikeByEvent").Return(0, nil)

		Service := eventService.NewEventService(
			eventRepositoryMock,
			userRepositoryMock,
			likeRepositoryMock,
		)
		actual, err := Service.Create(sampleRequest, int(sampleUser.ID), sampleFileRequest, storageProvider)

		expected := entities.EventResponse{}
		assert.Error(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("store-fail", func(t *testing.T) {
		sampleRequest := sampleRequestCentral
		sampleFileRequest := cover
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(sampleUser, nil)
		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("UploadFromRequest").Return("domain.com/images.jpg", nil)
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Store").Return(entities.Event{}, web.WebError{})
		eventRepositoryMock.Mock.On("Find").Return(entities.Event{}, nil)
		likeRepositoryMock := likeRepository.NewLikeRepositoryMock(&mock.Mock{})
		likeRepositoryMock.Mock.On("CountLikeByEvent").Return(0, nil)

		Service := eventService.NewEventService(
			eventRepositoryMock,
			userRepositoryMock,
			likeRepositoryMock,
		)
		actual, err := Service.Create(sampleRequest, int(sampleUser.ID), sampleFileRequest, storageProvider)

		expected := entities.EventResponse{}
		assert.Error(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("find-user-fail", func(t *testing.T) {
		sampleRequest := sampleRequestCentral
		sampleFileRequest := cover
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(entities.User{}, web.WebError{})
		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("UploadFromRequest").Return("domain.com/images.jpg", nil)
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Store").Return(entities.Event{}, web.WebError{})
		eventRepositoryMock.Mock.On("Find").Return(entities.Event{}, nil)
		likeRepositoryMock := likeRepository.NewLikeRepositoryMock(&mock.Mock{})
		likeRepositoryMock.Mock.On("CountLikeByEvent").Return(0, nil)

		Service := eventService.NewEventService(
			eventRepositoryMock,
			userRepositoryMock,
			likeRepositoryMock,
		)
		actual, err := Service.Create(sampleRequest, int(sampleUser.ID), sampleFileRequest, storageProvider)

		expected := entities.EventResponse{}
		assert.Error(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("find-event-fail", func(t *testing.T) {
		sampleRequest := sampleRequestCentral
		sampleFileRequest := cover
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(sampleUser, nil)
		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("UploadFromRequest").Return("domain.com/images.jpg", nil)
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Store").Return(entities.Event{}, nil)
		eventRepositoryMock.Mock.On("Find").Return(entities.Event{}, web.WebError{})
		likeRepositoryMock := likeRepository.NewLikeRepositoryMock(&mock.Mock{})
		likeRepositoryMock.Mock.On("CountLikeByEvent").Return(0, nil)

		Service := eventService.NewEventService(
			eventRepositoryMock,
			userRepositoryMock,
			likeRepositoryMock,
		)
		actual, err := Service.Create(sampleRequest, int(sampleUser.ID), sampleFileRequest, storageProvider)

		expected := entities.EventResponse{}
		assert.Error(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestUpdate(t *testing.T) {
	sampleCentral := eventRepository.EventCollection[0]
	sampleUser := userRepository.UserCollection[0]
	sampleRequestCentral := entities.EventRequest{}
	copier.Copy(&sampleRequestCentral, &sampleCentral)
	sampleRequestCentral.DatetimeEvent = "1999-12-12"
	cover := &multipart.FileHeader{
		Filename: "cover.jpg",
		Header: textproto.MIMEHeader{
			"Content-Disposition": []string{
				"form-data; name=\"cover\"; filename=\"cover.png\"",
			},
			"Content-Type": []string{
				"image/png",
			},
		},
		Size: 155 * 1024,
	}
	t.Run("success", func(t *testing.T) {
		sampleEvent := sampleCentral
		sampleRequest := sampleRequestCentral
		sampleFileRequest := cover
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Find").Return(sampleEvent, nil)
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(sampleUser, nil)
		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)
		storageProvider.Mock.On("UploadFromRequest").Return("example.com/images.png", nil)
		eventRepositoryMock.Mock.On("Update").Return(sampleEvent, nil)
		eventRepositoryMock.Mock.On("Find").Return(sampleEvent, nil)
		likeRepositoryMock := likeRepository.NewLikeRepositoryMock(&mock.Mock{})
		likeRepositoryMock.Mock.On("CountLikeByEvent").Return(0, nil)

		Service := eventService.NewEventService(
			eventRepositoryMock,
			userRepositoryMock,
			likeRepositoryMock,
		)
		actual, err := Service.Update(sampleRequest, 1, int(sampleUser.ID), sampleFileRequest, storageProvider)

		expected := entities.EventResponse{}
		copier.Copy(&expected, &sampleEvent)

		assert.Nil(t, err)
		assert.Equal(t, expected.ID, actual.ID)
	})
	t.Run("validation-fail", func(t *testing.T) {
		sampleEvent := sampleCentral
		sampleRequest := sampleRequestCentral
		sampleFileRequest := &multipart.FileHeader{
			Filename: "cover.jpg",
			Header: textproto.MIMEHeader{
				"Content-Disposition": []string{
					"form-data; name=\"cover\"; filename=\"cover.png\"",
				},
				"Content-Type": []string{
					"image/png",
				},
			},
			Size: 2024 * 2024,
		}

		sampleRequest.Title = ""
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(sampleUser, nil)
		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)
		storageProvider.Mock.On("UploadFromRequest").Return("example.com/images.png", nil)
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Update").Return(sampleEvent, web.WebError{})
		eventRepositoryMock.Mock.On("Find").Return(sampleEvent, nil)
		likeRepositoryMock := likeRepository.NewLikeRepositoryMock(&mock.Mock{})
		likeRepositoryMock.Mock.On("CountLikeByEvent").Return(0, nil)

		Service := eventService.NewEventService(
			eventRepositoryMock,
			userRepositoryMock,
			likeRepositoryMock,
		)
		actual, err := Service.Update(sampleRequest, 1, int(sampleUser.ID), sampleFileRequest, storageProvider)

		expected := entities.EventResponse{}
		copier.Copy(&expected, &sampleEvent)

		assert.Error(t, err)
		assert.Equal(t, entities.EventResponse{}, actual)
	})
	t.Run("invalid-datetime", func(t *testing.T) {
		sampleEvent := sampleCentral
		sampleRequest := sampleRequestCentral
		sampleFileRequest := cover

		sampleRequest.DatetimeEvent = "2022222222"
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(sampleUser, nil)
		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)
		storageProvider.Mock.On("UploadFromRequest").Return("example.com/images.png", nil)
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Update").Return(sampleEvent, nil)
		eventRepositoryMock.Mock.On("Find").Return(sampleEvent, nil)
		likeRepositoryMock := likeRepository.NewLikeRepositoryMock(&mock.Mock{})
		likeRepositoryMock.Mock.On("CountLikeByEvent").Return(0, nil)

		Service := eventService.NewEventService(
			eventRepositoryMock,
			userRepositoryMock,
			likeRepositoryMock,
		)
		actual, err := Service.Update(sampleRequest, 1, int(sampleUser.ID), sampleFileRequest, storageProvider)

		expected := entities.EventResponse{}
		copier.Copy(&expected, &sampleEvent)

		assert.Error(t, err)
		assert.Equal(t, entities.EventResponse{}, actual)
	})
	t.Run("upload-fail", func(t *testing.T) {
		sampleEvent := sampleCentral
		sampleRequest := sampleRequestCentral
		sampleFileRequest := cover
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(sampleUser, nil)
		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)
		storageProvider.Mock.On("UploadFromRequest").Return("", web.WebError{})
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Update").Return(sampleEvent, nil)
		eventRepositoryMock.Mock.On("Find").Return(sampleEvent, nil)
		likeRepositoryMock := likeRepository.NewLikeRepositoryMock(&mock.Mock{})
		likeRepositoryMock.Mock.On("CountLikeByEvent").Return(0, nil)

		Service := eventService.NewEventService(
			eventRepositoryMock,
			userRepositoryMock,
			likeRepositoryMock,
		)
		actual, err := Service.Update(sampleRequest, 1, int(sampleUser.ID), sampleFileRequest, storageProvider)

		expected := entities.EventResponse{}
		assert.Error(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("store-fail", func(t *testing.T) {
		sampleRequest := sampleRequestCentral
		sampleFileRequest := cover
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(sampleUser, nil)
		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)
		storageProvider.Mock.On("UploadFromRequest").Return("domain.com/images.jpg", nil)
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Update").Return(entities.Event{}, web.WebError{})
		eventRepositoryMock.Mock.On("Find").Return(entities.Event{}, nil)
		likeRepositoryMock := likeRepository.NewLikeRepositoryMock(&mock.Mock{})
		likeRepositoryMock.Mock.On("CountLikeByEvent").Return(0, nil)

		Service := eventService.NewEventService(
			eventRepositoryMock,
			userRepositoryMock,
			likeRepositoryMock,
		)
		actual, err := Service.Update(sampleRequest, 1, int(sampleUser.ID), sampleFileRequest, storageProvider)

		expected := entities.EventResponse{}
		assert.Error(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("find-user-fail", func(t *testing.T) {
		sampleRequest := sampleRequestCentral
		sampleFileRequest := cover
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(entities.User{}, web.WebError{})
		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)
		storageProvider.Mock.On("UploadFromRequest").Return("domain.com/images.jpg", nil)
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Update").Return(entities.Event{}, web.WebError{})
		eventRepositoryMock.Mock.On("Find").Return(entities.Event{}, nil)
		likeRepositoryMock := likeRepository.NewLikeRepositoryMock(&mock.Mock{})
		likeRepositoryMock.Mock.On("CountLikeByEvent").Return(0, nil)

		Service := eventService.NewEventService(
			eventRepositoryMock,
			userRepositoryMock,
			likeRepositoryMock,
		)
		actual, err := Service.Update(sampleRequest, 1, int(sampleUser.ID), sampleFileRequest, storageProvider)

		expected := entities.EventResponse{}
		assert.Error(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("find-event-fail", func(t *testing.T) {
		sampleRequest := sampleRequestCentral
		sampleFileRequest := cover
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(sampleUser, nil)
		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)
		storageProvider.Mock.On("UploadFromRequest").Return("domain.com/images.jpg", nil)
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Update").Return(entities.Event{}, nil)
		eventRepositoryMock.Mock.On("Find").Return(entities.Event{}, web.WebError{})
		likeRepositoryMock := likeRepository.NewLikeRepositoryMock(&mock.Mock{})
		likeRepositoryMock.Mock.On("CountLikeByEvent").Return(0, nil)

		Service := eventService.NewEventService(
			eventRepositoryMock,
			userRepositoryMock,
			likeRepositoryMock,
		)
		actual, err := Service.Update(sampleRequest, 1, int(sampleUser.ID), sampleFileRequest, storageProvider)

		expected := entities.EventResponse{}
		assert.Error(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		sampleEvent := eventRepository.EventCollection[0]
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Find").Return(sampleEvent, nil)

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)

		sampleUser := userRepository.UserCollection[0]
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(sampleUser, nil)

		likeRepositoryMock := likeRepository.NewLikeRepositoryMock(&mock.Mock{})
		likeRepositoryMock.Mock.On("CountLikeByEvent").Return(0, nil)

		eventRepositoryMock.Mock.On("Delete").Return(nil)

		Service := eventService.NewEventService(
			eventRepositoryMock,
			userRepositoryMock,
			likeRepositoryMock,
		)
		err := Service.Delete(int(sampleEvent.ID), int(sampleUser.ID), storageProvider)
		assert.Nil(t, err)
	})
	t.Run("repo-fail", func(t *testing.T) {
		sampleEvent := eventRepository.EventCollection[0]
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Find").Return(entities.Event{}, web.WebError{})

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)

		sampleUser := userRepository.UserCollection[0]
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(sampleUser, nil)

		likeRepositoryMock := likeRepository.NewLikeRepositoryMock(&mock.Mock{})
		likeRepositoryMock.Mock.On("CountLikeByEvent").Return(0, nil)
		eventRepositoryMock.Mock.On("Delete").Return(nil)

		Service := eventService.NewEventService(
			eventRepositoryMock,
			userRepositoryMock,
			likeRepositoryMock,
		)
		err := Service.Delete(int(sampleEvent.ID), int(sampleUser.ID), storageProvider)
		assert.Error(t, err)
	})
	t.Run("fail", func(t *testing.T) {
		sampleEvent := eventRepository.EventCollection[0]
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Find").Return(sampleEvent, nil)

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)

		sampleUser := userRepository.UserCollection[0]
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(sampleUser, nil)

		likeRepositoryMock := likeRepository.NewLikeRepositoryMock(&mock.Mock{})
		likeRepositoryMock.Mock.On("CountLikeByEvent").Return(0, nil)

		eventRepositoryMock.Mock.On("Delete").Return(web.WebError{})

		Service := eventService.NewEventService(
			eventRepositoryMock,
			userRepositoryMock,
			likeRepositoryMock,
		)
		err := Service.Delete(int(sampleEvent.ID), int(sampleUser.ID), storageProvider)
		assert.Error(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		sampleEvent := eventRepository.EventCollection[0]
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("Find").Return(sampleEvent, nil)

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)

		sampleUser := userRepository.UserCollection[0]
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(entities.User{}, web.WebError{})

		likeRepositoryMock := likeRepository.NewLikeRepositoryMock(&mock.Mock{})
		likeRepositoryMock.Mock.On("CountLikeByEvent").Return(0, nil)

		eventRepositoryMock.Mock.On("Delete").Return(nil)

		Service := eventService.NewEventService(
			eventRepositoryMock,
			userRepositoryMock,
			likeRepositoryMock,
		)
		err := Service.Delete(int(sampleEvent.ID), int(sampleUser.ID), storageProvider)
		assert.Error(t, err)
	})
}
