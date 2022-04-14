package handlers

import (
	"reflect"
	"strconv"
	"tupulung/config"
	"tupulung/entities"
	"tupulung/entities/web"
	categoryService "tupulung/services/category"

	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	categoryService *categoryService.CategoryService
}

func NewCategoryHandler(service *categoryService.CategoryService) CategoryHandler {
	return CategoryHandler{
		categoryService: service,
	}
}


/*
 * Category Handler - Index
 * -------------------------------
 * Mengambil list data category 
 * berdasarkan parameter query yang tersedia
 */
func (handler CategoryHandler) Index(c echo.Context) error {

	// Translate query param to map of filters
	filters := []map[string]string{}
	q := c.QueryParam("q") 
	if q != "" {
		filters = append(filters, map[string]string{
			"field": "title",
			"operator": "LIKE",
			"value": "%" + q + "%",
		})
	}

	// Sort parameter
	sorts := []map[string]interface{} {}
	sortPrice := c.QueryParam("sortPrice") 
	if sortPrice != "" {
		switch sortPrice {
		case "1":
			sorts = append(sorts, map[string]interface{} {
				"field": "title",
				"desc": true,
			})
		case "0":
			sorts = append(sorts, map[string]interface{} {
				"field": "title",
				"desc": false,
			})
		}
	}

	// Link hateoas
	links := map[string]string {"self": config.Get().App.BaseURL + "/api/categories"}

	// pagination param
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return c.JSON(400, web.ErrorResponse{
			Status: "ERROR",
			Code: 400,
			Error: "Limit parameter format is invalid",
			Links: links,
		})
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		links := map[string]string {"self": config.Get().App.BaseURL}
		return c.JSON(400, web.ErrorResponse{
			Status: "ERROR",
			Code: 400,
			Error: "Page parameter format is invalid",
			Links: links,
		})
	}

	// Get all categories
	categoriesRes, err := handler.categoryService.FindAll(limit, page, filters, sorts)
	if err != nil {
		if reflect.TypeOf(err).String() == "web.WebError" {
			webErr := err.(web.WebError)
			return c.JSON(webErr.Code, web.ErrorResponse{
				Status: "ERROR",
				Code: webErr.Code,
				Error: webErr.Error(),
				Links: links,
			})
		}
		return c.JSON(500, web.ErrorResponse{
			Status: "ERROR",
			Code: 500,
			Error: err.Error(),
			Links: links,
		})
	}
	
	// Get pagination data
	pagination, err := handler.categoryService.GetPagination(limit, page, filters)
	if err != nil {
		if reflect.TypeOf(err).String() == "web.WebError" {
			webErr := err.(web.WebError)
			return c.JSON(webErr.Code, web.ErrorResponse{
				Status: "ERROR",
				Code: webErr.Code,
				Error: webErr.Error(),
				Links: links,
			})
		}
		panic("not returning custom error")
	}
	
	// success response
	return c.JSON(200, web.SuccessListResponse{
		Status: "OK",
		Code: 200,
		Error: nil,
		Links: links,
		Data: categoriesRes,
		Pagination: pagination,
	})
}

/*
 * -------------------------------------------
 * Create category resource
 * -------------------------------------------
 */
func (handler CategoryHandler) Create(c echo.Context) error {
	// Populate form
	categoryReq := entities.CategoryRequest{}
	c.Bind(&categoryReq)
	
	// Define hateoas links
	links := map[string]string{ "self": config.Get().App.BaseURL + "/categories"}

	// Insert category
	categoryRes, err := handler.categoryService.Create(categoryReq)
	if err != nil {
		if reflect.TypeOf(err).String() == "web.WebError" {
			webErr := err.(web.WebError)
			return c.JSON(webErr.Code, web.ErrorResponse{
				Status: "ERROR",
				Code: webErr.Code,
				Error: webErr.Error(),
				Links: links,
			})
		} else if reflect.TypeOf(err).String() == "web.ValidationError" {
			valErr := err.(web.ValidationError)
			return c.JSON(valErr.Code, web.ValidationErrorResponse{
				Status: "ERROR",
				Code: valErr.Code,
				Error: valErr.Error(),
				Errors: valErr.Errors,
				Links: links,
			})
		}
	} 

	// response
	return c.JSON(200, web.SuccessResponse{
		Status: "OK",
		Code: 200,
		Error: nil,
		Links: links,
		Data: categoryRes,
	})
}


/*
 * -------------------------------------------
 * Update category resource
 * -------------------------------------------
 */
func (handler CategoryHandler) Update(c echo.Context) error {
	// Populate form
	categoryReq := entities.CategoryRequest{}
	c.Bind(&categoryReq)

	// Get param ID
	id, err := strconv.Atoi(c.Param("id")) 
	links := map[string]string{ "self": config.Get().App.BaseURL + "/categories/" + c.Param("id")}
	if err != nil {
		return c.JSON(400, web.ErrorResponse{
			Status: "ERROR",
			Code: 400,
			Error: "Parameter ID is invalid",
			Links: links,
		})
	}

	// Service call
	categoryRes, err := handler.categoryService.Update(categoryReq, id)
	if err != nil {
		if reflect.TypeOf(err).String() == "web.WebError" {
			webErr := err.(web.WebError)
			return c.JSON(webErr.Code, web.ErrorResponse{
				Status: "ERROR",
				Code: webErr.Code,
				Error: webErr.Error(),
				Links: links,
			})
		}
	}

	// response
	return c.JSON(200, web.SuccessResponse{
		Status: "OK",
		Code: 200,
		Error: nil,
		Links: links,
		Data: categoryRes,
	})
}


/*
 * -------------------------------------------
 * Delete category resource
 * -------------------------------------------
 */
func (handler CategoryHandler) Delete(c echo.Context) error {

	// Get params ID
	id, err := strconv.Atoi(c.Param("id")) 
	links := map[string]string{ "self": config.Get().App.BaseURL + "/categories/" + c.Param("id")}
	if err != nil {
		return c.JSON(400, web.ErrorResponse{
			Status: "ERROR",
			Code: 400,
			Error: "Parameter ID is invalid",
			Links: links,
		})
	}

	// call delete service
	err = handler.categoryService.Delete(id)
	if err != nil {
		if reflect.TypeOf(err).String() == "web.WebError" {
			webErr := err.(web.WebError)
			return c.JSON(webErr.Code, web.ErrorResponse{
				Status: "ERROR",
				Code: webErr.Code,
				Error: webErr.Error(),
				Links: links,
			})
		}
	}

	// response
	return c.JSON(200, web.SuccessResponse{
		Status: "OK",
		Code: 200,
		Error: nil,
		Links: links,
		Data: map[string]interface{} {
			"id": id,
		},
	})
}