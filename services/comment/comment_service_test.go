package comment_test

import (
	"testing"
	"tupulung/entities"
	"tupulung/entities/web"
	commentRepository "tupulung/repositories/comment"
	userRepository "tupulung/repositories/user"
	commentService "tupulung/services/comment"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFindAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		commentSample := commentRepository.CommentCollection
		commentRepositoryMock := commentRepository.NewCommentRepositoryMock(&mock.Mock{})
		commentRepositoryMock.Mock.On(
			"FindAll",
			0, 0,
			[]map[string]string{},
			[]map[string]interface{}{},
		).Return(commentSample, nil)

		service := commentService.NewCommentService(
			commentRepositoryMock,
			userRepository.NewUserRepositoryMock(&mock.Mock{}),
		)
		data, err := service.FindAll(0, 0, []map[string]string{}, []map[string]interface{}{})

		// Konversi expected data ke response
		commentRes := []entities.CommentResponse{}
		copier.Copy(&commentRes, &commentSample)

		assert.Nil(t, err)
		assert.Equal(t, commentRes, data)
	})
	t.Run("repo-fail", func(t *testing.T) {
		commentRepositoryMock := commentRepository.NewCommentRepositoryMock(&mock.Mock{})
		commentRepositoryMock.Mock.On(
			"FindAll",
			0, 0,
			[]map[string]string{},
			[]map[string]interface{}{},
		).Return([]entities.Comment{}, web.WebError{})

		service := commentService.NewCommentService(
			commentRepositoryMock,
			userRepository.NewUserRepositoryMock(&mock.Mock{}),
		)
		data, err := service.FindAll(0, 0, []map[string]string{}, []map[string]interface{}{})

		// Konversi expected data ke response
		commentRes := []entities.CommentResponse{}

		assert.Error(t, err)
		assert.Equal(t, commentRes, data)
	})
}

func TestGetPagination(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		commentRepositoryMock := commentRepository.NewCommentRepositoryMock(&mock.Mock{})
		commentRepositoryMock.Mock.On("CountAll", []map[string]string{}).Return(20, nil)

		Service := commentService.NewCommentService(
			commentRepositoryMock,
			userRepository.NewUserRepositoryMock(&mock.Mock{}),
		)
		actual, err := Service.GetPagination(1, 5, []map[string]string{})

		expected := web.Pagination{
			Page:       1,
			Limit:      5,
			TotalPages: int(4),
		}
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("repo-fail", func(t *testing.T) {
		commentRepositoryMock := commentRepository.NewCommentRepositoryMock(&mock.Mock{})
		commentRepositoryMock.Mock.On("CountAll", []map[string]string{}).Return(0, web.WebError{})

		Service := commentService.NewCommentService(
			commentRepositoryMock,
			userRepository.NewUserRepositoryMock(&mock.Mock{}),
		)
		actual, err := Service.GetPagination(5, 1, []map[string]string{})
		assert.Error(t, err)
		assert.Equal(t, web.Pagination{}, actual)
	})
	t.Run("limit-zero", func(t *testing.T) {
		commentRepositoryMock := commentRepository.NewCommentRepositoryMock(&mock.Mock{})
		commentRepositoryMock.Mock.On("CountAll", []map[string]string{}).Return(20, nil)

		Service := commentService.NewCommentService(
			commentRepositoryMock,
			userRepository.NewUserRepositoryMock(&mock.Mock{}),
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
		commentRepositoryMock := commentRepository.NewCommentRepositoryMock(&mock.Mock{})
		commentRepositoryMock.Mock.On("CountAll", []map[string]string{}).Return(0, nil)

		Service := commentService.NewCommentService(
			commentRepositoryMock,
			userRepository.NewUserRepositoryMock(&mock.Mock{}),
		)
		actual, err := Service.GetPagination(0, 1, []map[string]string{})

		expected := web.Pagination{
			Page:       0,
			Limit:      1,
			TotalPages: int(1),
		}
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("added-page-on-active-module", func(t *testing.T) {
		commentRepositoryMock := commentRepository.NewCommentRepositoryMock(&mock.Mock{})
		commentRepositoryMock.Mock.On("CountAll", []map[string]string{}).Return(20, nil)

		Service := commentService.NewCommentService(
			commentRepositoryMock,
			userRepository.NewUserRepositoryMock(&mock.Mock{}),
		)
		actual, err := Service.GetPagination(1, 5, []map[string]string{})

		expected := web.Pagination{
			Page:       1,
			Limit:      5,
			TotalPages: int(4),
		}
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestCreate(t *testing.T) {

	sampleCommentCentral := commentRepository.CommentCollection[0]
	sampleRequestCentral := entities.CommentRequest{}
	copier.Copy(&sampleRequestCentral, &sampleCommentCentral)

	t.Run("success", func(t *testing.T) {
		userSample := userRepository.UserCollection[0]
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(userSample, nil)
		sampleComment := sampleCommentCentral
		sampleRequest := sampleRequestCentral

		commentRepositoryMock := commentRepository.NewCommentRepositoryMock(&mock.Mock{})
		commentRepositoryMock.Mock.On("Store").Return(sampleComment, nil)

		Service := commentService.NewCommentService(
			commentRepositoryMock,
			userRepositoryMock,
		)
		actual, err := Service.Create(sampleRequest, 1, int(userSample.ID))

		expected := entities.CommentResponse{}
		copier.Copy(&expected, &sampleComment)

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("repo-fail", func(t *testing.T) {
		userSample := userRepository.UserCollection[0]
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(userSample, nil)
		sampleRequest := sampleRequestCentral

		commentRepositoryMock := commentRepository.NewCommentRepositoryMock(&mock.Mock{})
		commentRepositoryMock.Mock.On("Store").Return(entities.Comment{}, web.WebError{})

		Service := commentService.NewCommentService(
			commentRepositoryMock,
			userRepositoryMock,
		)
		actual, err := Service.Create(sampleRequest, 1, int(userSample.ID))

		expected := entities.CommentResponse{}

		assert.Error(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("find-user-fail", func(t *testing.T) {
		userSample := userRepository.UserCollection[0]
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(entities.User{}, web.WebError{})
		sampleRequest := sampleRequestCentral

		commentRepositoryMock := commentRepository.NewCommentRepositoryMock(&mock.Mock{})
		commentRepositoryMock.Mock.On("Store").Return(entities.Comment{}, nil)

		Service := commentService.NewCommentService(
			commentRepositoryMock,
			userRepositoryMock,
		)
		actual, err := Service.Create(sampleRequest, 1, int(userSample.ID))

		expected := entities.CommentResponse{}

		assert.Error(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("validation-fail", func(t *testing.T) {
		userSample := userRepository.UserCollection[0]
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(userSample, nil)
		sampleComment := sampleCommentCentral
		sampleRequest := sampleRequestCentral

		sampleRequest.Comment = ""
		commentRepositoryMock := commentRepository.NewCommentRepositoryMock(&mock.Mock{})
		commentRepositoryMock.Mock.On("Store").Return(sampleComment, nil)

		Service := commentService.NewCommentService(
			commentRepositoryMock,
			userRepositoryMock,
		)
		actual, err := Service.Create(sampleRequest, 1, int(userSample.ID))

		expected := entities.CommentResponse{}
		copier.Copy(&expected, &sampleComment)

		assert.Error(t, err)
		assert.Equal(t, entities.CommentResponse{}, actual)
	})

}

func TestUpdate(t *testing.T) {
	sampleCommentCentral := commentRepository.CommentCollection[0]
	sampleRequestCentral := entities.CommentRequest{}
	copier.Copy(&sampleRequestCentral, &sampleCommentCentral)

	t.Run("success", func(t *testing.T) {
		userSample := userRepository.UserCollection[0]
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(userSample, nil)
		sampleRequest := sampleRequestCentral
		sampleComment := sampleCommentCentral

		commentRepositoryMock := commentRepository.NewCommentRepositoryMock(&mock.Mock{})
		commentRepositoryMock.Mock.On("Find").Return(sampleComment, nil)

		commentOutput := sampleComment
		copier.CopyWithOption(&commentOutput, &sampleRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		commentRepositoryMock.Mock.On("Update").Return(commentOutput, nil)

		Service := commentService.NewCommentService(
			commentRepositoryMock,
			userRepositoryMock,
		)
		actual, err := Service.Update(sampleRequest, int(sampleComment.ID), int(userSample.ID))
		expected := entities.CommentResponse{}
		copier.Copy(&expected, &commentOutput)

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("repo-fail", func(t *testing.T) {
		userSample := userRepository.UserCollection[0]
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(userSample, nil)
		sampleRequest := sampleRequestCentral
		sampleComment := sampleCommentCentral

		commentRepositoryMock := commentRepository.NewCommentRepositoryMock(&mock.Mock{})
		commentRepositoryMock.Mock.On("Find").Return(entities.Comment{}, web.WebError{})

		commentOutput := sampleComment
		copier.CopyWithOption(&commentOutput, &sampleRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		commentRepositoryMock.Mock.On("Update").Return(commentOutput, nil)

		Service := commentService.NewCommentService(
			commentRepositoryMock,
			userRepositoryMock,
		)
		actual, err := Service.Update(sampleRequest, int(sampleComment.ID), int(userSample.ID))
		expected := entities.CommentResponse{}

		assert.Error(t, err)
		assert.Equal(t, expected, actual)
	})

}

func TestDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userSample := userRepository.UserCollection[0]
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(userSample, nil)
		sampleComment := commentRepository.CommentCollection[0]
		commentRepositoryMock := commentRepository.NewCommentRepositoryMock(&mock.Mock{})
		commentRepositoryMock.Mock.On("Find").Return(sampleComment, nil)

		commentRepositoryMock.Mock.On("Delete").Return(nil)

		Service := commentService.NewCommentService(
			commentRepositoryMock,
			userRepositoryMock,
		)
		err := Service.Delete(int(sampleComment.ID), int(userSample.ID))
		assert.Nil(t, err)
	})
	t.Run("repo-fail", func(t *testing.T) {
		userSample := userRepository.UserCollection[0]
		userRepositoryMock := userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(userSample, nil)
		sampleComment := commentRepository.CommentCollection[0]
		commentRepositoryMock := commentRepository.NewCommentRepositoryMock(&mock.Mock{})
		commentRepositoryMock.Mock.On("Find").Return(entities.Comment{}, web.WebError{})

		commentRepositoryMock.Mock.On("Delete").Return(nil)

		Service := commentService.NewCommentService(
			commentRepositoryMock,
			userRepositoryMock,
		)
		err := Service.Delete(int(sampleComment.ID), int(userSample.ID))
		assert.Error(t, err)
	})
}
