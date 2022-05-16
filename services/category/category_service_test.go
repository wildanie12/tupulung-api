package category_test

import (
	"testing"
	"tupulung/entities"
	web "tupulung/entities/web"
	categoryRepostory "tupulung/repositories/category"
	categoryService "tupulung/services/category"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFindAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		categorySample := categoryRepostory.CategoryCollection
		categoryRepositoryMock := categoryRepostory.NewCategoryRepositoryMock(&mock.Mock{})
		categoryRepositoryMock.Mock.On(
			"FindAll",
			0, 0,
			[]map[string]string{},
			[]map[string]interface{}{},
		).Return(categorySample, nil)

		service := categoryService.NewCategoryService(categoryRepositoryMock)
		data, err := service.FindAll(0, 0, []map[string]string{}, []map[string]interface{}{})

		// Konversi expected data ke response
		categoryRes := []entities.CategoryResponse{}
		copier.Copy(&categoryRes, &categorySample)

		assert.Nil(t, err)
		assert.Equal(t, categoryRes, data)
	})
	t.Run("success", func(t *testing.T) {
		categoryRepositoryMock := categoryRepostory.NewCategoryRepositoryMock(&mock.Mock{})
		categoryRepositoryMock.Mock.On(
			"FindAll",
			0, 0,
			[]map[string]string{},
			[]map[string]interface{}{},
		).Return([]entities.Category{}, web.WebError{})

		service := categoryService.NewCategoryService(categoryRepositoryMock)
		data, err := service.FindAll(0, 0, []map[string]string{}, []map[string]interface{}{})

		// Konversi expected data ke response
		categoryRes := []entities.CategoryResponse{}

		assert.Error(t, err)
		assert.Equal(t, categoryRes, data)
	})
}

func TestGetPagination(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		categoryRepositoryMock := categoryRepostory.NewCategoryRepositoryMock(&mock.Mock{})
		categoryRepositoryMock.Mock.On("CountAll", []map[string]string{}).Return(20, nil)

		Service := categoryService.NewCategoryService(
			categoryRepositoryMock,
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
		categoryRepositoryMock := categoryRepostory.NewCategoryRepositoryMock(&mock.Mock{})
		categoryRepositoryMock.Mock.On("CountAll", []map[string]string{}).Return(0, web.WebError{})

		Service := categoryService.NewCategoryService(
			categoryRepositoryMock,
		)
		actual, err := Service.GetPagination(5, 1, []map[string]string{})
		assert.Error(t, err)
		assert.Equal(t, web.Pagination{}, actual)
	})
	t.Run("limit-zero", func(t *testing.T) {
		categoryRepositoryMock := categoryRepostory.NewCategoryRepositoryMock(&mock.Mock{})
		categoryRepositoryMock.Mock.On("CountAll", []map[string]string{}).Return(20, nil)

		Service := categoryService.NewCategoryService(
			categoryRepositoryMock,
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
	t.Run("added-page-on-active-module", func(t *testing.T) {
		categoryRepositoryMock := categoryRepostory.NewCategoryRepositoryMock(&mock.Mock{})
		categoryRepositoryMock.Mock.On("CountAll", []map[string]string{}).Return(20, nil)

		Service := categoryService.NewCategoryService(
			categoryRepositoryMock,
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

func TestFind(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		categoryOutput := categoryRepostory.CategoryCollection[0]
		categoryRepositoryMock := categoryRepostory.NewCategoryRepositoryMock(&mock.Mock{})
		categoryRepositoryMock.Mock.On("Find", 1).Return(categoryOutput, nil)

		orderService := categoryService.NewCategoryService(
			categoryRepositoryMock,
		)
		actual, err := orderService.Find(int(categoryOutput.ID))

		// convert to response
		expected := entities.CategoryResponse{}
		copier.Copy(&expected, &categoryOutput)
		expected.ID = categoryOutput.ID // fix: overlap destinationID vs orderID

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("error", func(t *testing.T) {
		categoryRepositoryMock := categoryRepostory.NewCategoryRepositoryMock(&mock.Mock{})
		categoryRepositoryMock.Mock.On("Find", 1).Return(entities.Category{}, web.WebError{Code: 500, Message: "Error"})

		orderService := categoryService.NewCategoryService(
			categoryRepositoryMock,
		)
		actual, err := orderService.Find(1)
		assert.Error(t, err)
		assert.Equal(t, entities.CategoryResponse{}, actual)
	})
}

func TestCreate(t *testing.T) {
	sampleCategoryCentral := categoryRepostory.CategoryCollection[0]
	sampleRequestCentral := entities.CategoryRequest{}
	copier.Copy(&sampleRequestCentral, &sampleCategoryCentral)

	t.Run("success", func(t *testing.T) {
		sampleCategory := sampleCategoryCentral
		sampleRequest := sampleRequestCentral

		categoryRepositoryMock := categoryRepostory.NewCategoryRepositoryMock(&mock.Mock{})
		categoryRepositoryMock.Mock.On("Store").Return(sampleCategory, nil)

		Service := categoryService.NewCategoryService(
			categoryRepositoryMock,
		)
		actual, err := Service.Create(sampleRequest)

		expected := entities.CategoryResponse{}
		copier.Copy(&expected, &sampleCategory)

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("validation-fail", func(t *testing.T) {
		sampleCategory := sampleCategoryCentral
		sampleRequest := sampleRequestCentral

		sampleRequest.Title = ""
		categoryRepositoryMock := categoryRepostory.NewCategoryRepositoryMock(&mock.Mock{})
		categoryRepositoryMock.Mock.On("Store").Return(sampleCategory, nil)

		Service := categoryService.NewCategoryService(
			categoryRepositoryMock,
		)
		actual, err := Service.Create(sampleRequest)

		expected := entities.CategoryResponse{}
		copier.Copy(&expected, &sampleCategory)

		assert.Error(t, err)
		assert.Equal(t, entities.CategoryResponse{}, actual)
	})

}

func TestUpdate(t *testing.T) {
	sampleCategoryCentral := categoryRepostory.CategoryCollection[0]
	sampleRequestCentral := entities.CategoryRequest{}
	copier.Copy(&sampleRequestCentral, &sampleCategoryCentral)

	t.Run("success", func(t *testing.T) {
		sampleRequest := sampleRequestCentral
		sampleCategory := sampleCategoryCentral

		categoryRepositoryMock := categoryRepostory.NewCategoryRepositoryMock(&mock.Mock{})
		categoryRepositoryMock.Mock.On("Find", 1).Return(sampleCategory, nil)

		categoryOutput := sampleCategory
		copier.CopyWithOption(&categoryOutput, &sampleRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		categoryRepositoryMock.Mock.On("Update").Return(categoryOutput, nil)

		Service := categoryService.NewCategoryService(
			categoryRepositoryMock,
		)
		actual, err := Service.Update(sampleRequest, int(sampleCategory.ID))
		expected := entities.CategoryResponse{}
		copier.Copy(&expected, &categoryOutput)

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		sampleCategory := categoryRepostory.CategoryCollection[0]
		categoryRepositoryMock := categoryRepostory.NewCategoryRepositoryMock(&mock.Mock{})
		categoryRepositoryMock.Mock.On("Find", 1).Return(sampleCategory, nil)

		categoryRepositoryMock.Mock.On("Delete", 1).Return(nil)

		categoryService := categoryService.NewCategoryService(
			categoryRepositoryMock,
		)
		err := categoryService.Delete(int(sampleCategory.ID))
		assert.Nil(t, err)
	})
	t.Run("repo-fail", func(t *testing.T) {
		sampleCategory := categoryRepostory.CategoryCollection[0]
		categoryRepositoryMock := categoryRepostory.NewCategoryRepositoryMock(&mock.Mock{})
		categoryRepositoryMock.Mock.On("Find", 1).Return(entities.Category{}, web.WebError{})

		categoryRepositoryMock.Mock.On("Delete", 1).Return(nil)

		categoryService := categoryService.NewCategoryService(
			categoryRepositoryMock,
		)
		err := categoryService.Delete(int(sampleCategory.ID))
		assert.Error(t, err)
	})
}
