package user_test

import (
	"mime/multipart"
	"net/textproto"
	"testing"

	"tupulung/entities"
	"tupulung/entities/web"
	eventRepository "tupulung/repositories/event"
	userRepository "tupulung/repositories/user"
	userService "tupulung/services/user"
	_storageProvider "tupulung/utilities/storage"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFind(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		userSample := userRepository.UserCollection[0]
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(userSample, nil)

		Service := userService.NewUserService(
			userRepositoryMock,
			eventRepositoryMock,
		)
		_, err := Service.Find(int(userSample.ID))

		assert.Nil(t, err)
	})
}
func TestGetJoinedEvent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		eventSample := eventRepository.EventCollection
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		userSample := userRepository.UserCollection[0]
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("GetJoinedEvents").Return(eventSample, nil)

		Service := userService.NewUserService(
			userRepositoryMock,
			eventRepositoryMock,
		)
		actual, err := Service.GetJoinedEvents(int(userSample.ID))
		expected := []entities.EventResponse{}
		copier.Copy(&expected, &eventSample)

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("repo-fail", func(t *testing.T) {
		eventSample := eventRepository.EventCollection
		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		userSample := userRepository.UserCollection[0]
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("GetJoinedEvents").Return([]entities.Event{}, web.WebError{})

		Service := userService.NewUserService(
			userRepositoryMock,
			eventRepositoryMock,
		)
		actual, err := Service.GetJoinedEvents(int(userSample.ID))
		expected := []entities.EventResponse{}
		copier.Copy(&expected, &eventSample)

		assert.Error(t, err)
		assert.Equal(t, []entities.EventResponse{}, actual)
	})
}

func TestCreate(t *testing.T) {
	sampleCentral := userRepository.UserCollection[0]
	sampleRequestCentral := entities.UserRequest{}
	copier.Copy(&sampleRequestCentral, &sampleCentral)
	sampleRequestCentral.DOB = "1999-12-12"
	avatar := &multipart.FileHeader{
		Filename: "avatar.jpg",
		Header: textproto.MIMEHeader{
			"Content-Disposition": []string{
				"form-data; name=\"avatar\"; filename=\"avatar.png\"",
			},
			"Content-Type": []string{
				"image/png",
			},
		},
		Size: 155 * 1024,
	}
	t.Run("success", func(t *testing.T) {
		sampleUser := sampleCentral
		sampleRequest := sampleRequestCentral
		sampleFileRequest := avatar

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("UploadFromRequest").Return("example.com/images.png", nil)
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Store").Return(sampleUser, nil)

		Service := userService.NewUserService(
			userRepositoryMock,
			eventRepository.NewEventRepositoryMock(&mock.Mock{}),
		)
		actual, err := Service.Create(sampleRequest, sampleFileRequest, storageProvider)

		expected := entities.UserResponse{}
		copier.Copy(&expected, &sampleUser)

		assert.Nil(t, err)
		assert.NotEqual(t, "", actual.Token)
		assert.Equal(t, expected, actual.User)
	})
	t.Run("validation-fail", func(t *testing.T) {
		sampleUser := sampleCentral
		sampleRequest := sampleRequestCentral
		sampleFileRequest := avatar

		sampleRequest.Name = ""
		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("UploadFromRequest").Return("example.com/images.png", nil)
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Store").Return(sampleUser, nil)

		Service := userService.NewUserService(
			userRepositoryMock,
			eventRepository.NewEventRepositoryMock(&mock.Mock{}),
		)
		actual, err := Service.Create(sampleRequest, sampleFileRequest, storageProvider)

		expected := entities.UserResponse{}
		copier.Copy(&expected, &sampleUser)

		assert.Error(t, err)
		assert.Equal(t, entities.AuthResponse{}, actual)
	})
	t.Run("invalid-dob", func(t *testing.T) {
		sampleUser := sampleCentral
		sampleRequest := sampleRequestCentral
		sampleFileRequest := avatar

		sampleRequest.DOB = "2022222222"
		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("UploadFromRequest").Return("example.com/images.png", nil)
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Store").Return(sampleUser, nil)

		Service := userService.NewUserService(
			userRepositoryMock,
			eventRepository.NewEventRepositoryMock(&mock.Mock{}),
		)
		actual, err := Service.Create(sampleRequest, sampleFileRequest, storageProvider)

		expected := entities.UserResponse{}
		copier.Copy(&expected, &sampleUser)

		assert.Error(t, err)
		assert.Equal(t, entities.AuthResponse{}, actual)
	})
	t.Run("upload-fail", func(t *testing.T) {
		sampleUser := sampleCentral
		sampleRequest := sampleRequestCentral
		sampleFileRequest := avatar

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("UploadFromRequest").Return("", web.WebError{})
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Store").Return(sampleUser, nil)

		Service := userService.NewUserService(
			userRepositoryMock,
			eventRepository.NewEventRepositoryMock(&mock.Mock{}),
		)
		actual, err := Service.Create(sampleRequest, sampleFileRequest, storageProvider)

		expected := entities.AuthResponse{}
		assert.Error(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("store-fail", func(t *testing.T) {
		sampleRequest := sampleRequestCentral
		sampleFileRequest := avatar

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("UploadFromRequest").Return("domain.com/images.jpg", nil)
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Store").Return(entities.User{}, web.WebError{})

		Service := userService.NewUserService(
			userRepositoryMock,
			eventRepository.NewEventRepositoryMock(&mock.Mock{}),
		)
		actual, err := Service.Create(sampleRequest, sampleFileRequest, storageProvider)

		expected := entities.AuthResponse{}
		assert.Error(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestUpdate(t *testing.T) {
	sampleUserCentral := userRepository.UserCollection[0]
	sampleRequestCentral := entities.UserRequest{}
	copier.Copy(&sampleRequestCentral, &sampleUserCentral)
	sampleRequestCentral.DOB = "1999-12-12"
	avatar := &multipart.FileHeader{
		Filename: "avatar.jpg",
		Header: textproto.MIMEHeader{
			"Content-Disposition": []string{
				"form-data; name=\"avatar\"; filename=\"avatar.png\"",
			},
			"Content-Type": []string{
				"image/png",
			},
		},
		Size: 155 * 1024,
	}
	t.Run("success", func(t *testing.T) {
		sampleRequest := sampleRequestCentral
		sampleFileRequest := avatar
		sampleUser := sampleUserCentral

		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(sampleUser, nil)

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)
		storageProvider.Mock.On("UploadFromRequest").Return("domain.com/images.jpg", nil)

		userOutput := sampleUser
		copier.CopyWithOption(&userOutput, &sampleRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		userRepositoryMock.Mock.On("Update").Return(userOutput, nil)

		Service := userService.NewUserService(
			userRepositoryMock,
			eventRepository.NewEventRepositoryMock(&mock.Mock{}),
		)
		actual, err := Service.Update(sampleRequest, int(sampleUser.ID), sampleFileRequest, storageProvider)
		expected := entities.UserResponse{}
		copier.Copy(&expected, &userOutput)

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("validation-fail", func(t *testing.T) {
		sampleRequest := sampleRequestCentral
		sampleFileRequest := avatar
		sampleUser := sampleUserCentral
		avatar.Size = 1024 * 2048

		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(sampleUser, nil)

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)
		storageProvider.Mock.On("UploadFromRequest").Return("domain.com/images.jpg", nil)

		userOutput := sampleUser
		copier.CopyWithOption(&userOutput, &sampleRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		userRepositoryMock.Mock.On("Update").Return(userOutput, nil)

		Service := userService.NewUserService(
			userRepositoryMock,
			eventRepository.NewEventRepositoryMock(&mock.Mock{}),
		)
		actual, err := Service.Update(sampleRequest, int(sampleUser.ID), sampleFileRequest, storageProvider)
		assert.Error(t, err)
		assert.Equal(t, entities.UserResponse{}, actual)
	})
	t.Run("find-fail", func(t *testing.T) {
		sampleRequest := sampleRequestCentral
		sampleFileRequest := &multipart.FileHeader{
			Filename: "avatar.jpg",
			Header: textproto.MIMEHeader{
				"Content-Disposition": []string{
					"form-data; name=\"avatar\"; filename=\"avatar.png\"",
				},
				"Content-Type": []string{
					"image/png",
				},
			},
		}
		sampleUser := sampleUserCentral

		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(entities.User{}, web.WebError{})

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)
		storageProvider.Mock.On("UploadFromRequest").Return("domain.com/images.jpg", nil)

		userOutput := sampleUser
		copier.CopyWithOption(&userOutput, &sampleRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		userRepositoryMock.Mock.On("Update").Return(userOutput, nil)

		Service := userService.NewUserService(
			userRepositoryMock,
			eventRepository.NewEventRepositoryMock(&mock.Mock{}),
		)
		actual, err := Service.Update(sampleRequest, int(sampleUser.ID), sampleFileRequest, storageProvider)
		assert.Error(t, err)
		assert.Equal(t, entities.UserResponse{}, actual)
	})
	t.Run("upload-fail", func(t *testing.T) {
		sampleRequest := sampleRequestCentral
		sampleFileRequest := &multipart.FileHeader{
			Filename: "avatar.jpg",
			Header: textproto.MIMEHeader{
				"Content-Disposition": []string{
					"form-data; name=\"avatar\"; filename=\"avatar.png\"",
				},
				"Content-Type": []string{
					"image/png",
				},
			},
		}
		sampleUser := sampleUserCentral

		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(sampleUser, nil)

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)
		storageProvider.Mock.On("UploadFromRequest").Return("", web.WebError{})

		userOutput := sampleUser
		copier.CopyWithOption(&userOutput, &sampleRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		userRepositoryMock.Mock.On("Update").Return(userOutput, nil)

		Service := userService.NewUserService(
			userRepositoryMock,
			eventRepository.NewEventRepositoryMock(&mock.Mock{}),
		)
		actual, err := Service.Update(sampleRequest, int(sampleUser.ID), sampleFileRequest, storageProvider)
		assert.Error(t, err)
		assert.Equal(t, entities.UserResponse{}, actual)
	})
}

func TestDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		sampleCustomer := userRepository.UserCollection[0]
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(sampleCustomer, nil)

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)

		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("DeleteBatch").Return(nil)

		userRepositoryMock.Mock.On("Delete").Return(nil)

		Service := userService.NewUserService(
			userRepositoryMock,
			eventRepositoryMock,
		)
		err := Service.Delete(int(sampleCustomer.ID), storageProvider)
		assert.Nil(t, err)
	})
	t.Run("repo-fail", func(t *testing.T) {
		sampleCustomer := userRepository.UserCollection[0]
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(entities.User{}, web.WebError{})

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)

		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("DeleteBatch").Return(nil)

		userRepositoryMock.Mock.On("Delete").Return(nil)

		Service := userService.NewUserService(
			userRepositoryMock,
			eventRepositoryMock,
		)
		err := Service.Delete(int(sampleCustomer.ID), storageProvider)
		assert.Error(t, err)
	})
	t.Run("fail", func(t *testing.T) {
		sampleCustomer := userRepository.UserCollection[0]
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(sampleCustomer, nil)

		storageProvider := _storageProvider.NewStorageMock(&mock.Mock{})
		storageProvider.Mock.On("Delete").Return(nil)

		eventRepositoryMock := eventRepository.NewEventRepositoryMock(&mock.Mock{})
		eventRepositoryMock.Mock.On("DeleteBatch").Return(nil)

		userRepositoryMock.Mock.On("Delete").Return(web.WebError{})

		Service := userService.NewUserService(
			userRepositoryMock,
			eventRepositoryMock,
		)
		err := Service.Delete(int(sampleCustomer.ID), storageProvider)
		assert.Error(t, err)
	})
}
