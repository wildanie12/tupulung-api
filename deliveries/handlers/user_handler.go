package handlers

import (
	"net/http"
	"reflect"
	"strconv"
	"tupulung/config"
	"tupulung/entities"
	"tupulung/entities/web"
	userService "tupulung/services/user"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService *userService.UserService
}

func NewUserHandler(service *userService.UserService) *UserHandler {
	return &UserHandler{
		userService: service,
	}
}


/*
 * User Handler - Create
 * -------------------------------
 * Registrasi User kedalam sistem dan
 * mengembalikan token 
 */
func (handler UserHandler) Create(c echo.Context) error {

	// Bind request ke user request
	userReq := entities.UserRequest{}
	c.Bind(&userReq)
	
	// Define links (hateoas)
	links := map[string]string{ "self": config.Get().App.BaseURL + "/api/users"}

	// Read file avatar
	avatar, _ := c.FormFile("avatar")

	// registrasi user via call user service
	userRes, err := handler.userService.Create(userReq, avatar)
	if err != nil {

		// return error response khusus jika err termasuk webError / ValidationError
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

		// return error 500 jika bukan webError
		return c.JSON(http.StatusInternalServerError, web.ErrorResponse{
			Status: "ERROR",
			Code: http.StatusInternalServerError,
			Error: err.Error(),
			Links: links,
		})
	}

	// response
	return c.JSON(http.StatusCreated, web.SuccessResponse{
		Status: "OK",
		Code: http.StatusCreated,
		Error: nil,
		Links: links,
		Data: userRes,
	})
}



/*
 * User Handler - Show
 * -------------------------------
 * Mendapatkan data user tunggal
 * berdasarkan ID di parameter path 
 */
func (handler UserHandler) Show(c echo.Context) error {

	// Get parameter ID and Set links hateoas
	id, err := strconv.Atoi(c.Param("id"))
	links := map[string]string{ "self": config.Get().App.BaseURL + "/api/users/" + c.Param("id") }
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Code: http.StatusBadGateway,
			Status: "ERROR",
			Error: "Invalid parameter",
			Links: links,
		})
	}

	// Get userdata via service call
	user, err := handler.userService.Find(id)
	if err != nil {
		if reflect.TypeOf(err).String() == "web.WebError" {
			webErr := err.(web.WebError)
			return c.JSON(webErr.Code, web.ErrorResponse{
				Code: webErr.Code,
				Status: "ERROR",
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
		Data: user,
	})
}

/*
 * User Handler - Update
 * -------------------------------
 * Edit profile user berdasarkan ID
 * user hanya dapat merubah data usernya sendiri
 */
func (handler UserHandler) Update(c echo.Context) error {

	// Bind request to user request
	userReq := entities.UserRequest{}
	c.Bind(&userReq)

	// Get parameter ID and set links hateoas
	id, err := strconv.Atoi(c.Param("id")) 
	links := map[string]string{ "self": config.Get().App.BaseURL + "/users/" + c.Param("id")}
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Code: http.StatusBadGateway,
			Status: "ERROR",
			Error: "Invalid parameter",
			Links: links,
		})
	}

	// Get token
	token := c.Get("user")
	
	// avatar
	avatar, _ := c.FormFile("avatar")

	// Update via user service call
	userRes, err := handler.userService.Update(userReq, id, avatar, token)
	if err != nil {
		if reflect.TypeOf(err).String() == "web.WebError" {
			webErr := err.(web.WebError)
			return c.JSON(webErr.Code, web.ErrorResponse{
				Code: webErr.Code,
				Status: "ERROR",
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
		Data: map[string]interface{} {
			"id": userRes.ID,
		},
	})
}


/*
 * User Handler - Delete
 * -------------------------------
 * Delete User dari sistem 
 * Hanya usernya sendiri yang dapat menghapus
 */
func (handler UserHandler) Delete(c echo.Context) error {

	// Get params ID and set links hateoas
	id, err := strconv.Atoi(c.Param("id")) 
	links := map[string]string{ "self": config.Get().App.BaseURL + "/users/" + c.Param("id")}
	if err != nil {
		return c.JSON(http.StatusBadRequest, web.ErrorResponse{
			Code: http.StatusBadGateway,
			Status: "ERROR",
			Error: "Invalid parameter",
			Links: links,
		})
	}

	// Get token
	token := c.Get("user")

	// call delete service
	err = handler.userService.Delete(id,token)
	if err != nil {
		if reflect.TypeOf(err).String() == "web.WebError" {
			webErr := err.(web.WebError)
			return c.JSON(webErr.Code, web.ErrorResponse{
				Code: webErr.Code,
				Status: "ERROR",
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