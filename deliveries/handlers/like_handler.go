package handlers

import (
	"net/http"
	"reflect"
	"strconv"
	"tupulung/config"
	"tupulung/deliveries/helpers"
	"tupulung/deliveries/middleware"
	"tupulung/entities/web"
	likeService "tupulung/services/like"

	"github.com/labstack/echo/v4"
)

type LikeHandler struct {
	likeService *likeService.LikeService
}

func NewLikeHandler(service *likeService.LikeService) *LikeHandler {
	return &LikeHandler{
		likeService: service,
	}
}

func (handler LikeHandler) Append(c echo.Context) error {

	token := c.Get("user")
	ID, tx := middleware.ReadToken(token)

	eventID, err := strconv.Atoi(c.Param("id"))
	links := map[string]string{"self": config.Get().App.BaseURL + "/events/" + c.Param("id")}
	if tx != nil {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Code:   http.StatusUnauthorized,
			Status: "ERROR",
			Error:  "unauthorized",
			Links:  links,
		})
	}
	if err != nil {
		return c.JSON(400, helpers.MakeErrorResponse("ERROR", 400, err.Error(), links))
	}

	tx = handler.likeService.Append(ID, eventID)

	if tx != nil {
		if reflect.TypeOf(tx).String() == "web.WebError" {
			webErr := tx.(web.WebError)
			return c.JSON(webErr.Code, helpers.MakeErrorResponse("ERROR", webErr.Code, webErr.Error(), links))
		}
		return c.JSON(500, helpers.MakeErrorResponse("ERROR", 500, tx.Error(), links))

	}

	return c.JSON(200, web.SuccessResponse{
		Status: "OK",
		Code:   200,
		Error:  nil,
		Links:  links,
		Data:   "Success like this event",
	})
}

func (handler LikeHandler) Delete(c echo.Context) error {

	token := c.Get("user")
	ID, tx := middleware.ReadToken(token)

	eventID, err := strconv.Atoi(c.Param("id"))
	links := map[string]string{"self": config.Get().App.BaseURL + "/events/" + c.Param("id")}
	if tx != nil {
		return c.JSON(http.StatusUnauthorized, web.ErrorResponse{
			Code:   http.StatusUnauthorized,
			Status: "ERROR",
			Error:  "unauthorized",
			Links:  links,
		})
	}
	if err != nil {
		return c.JSON(400, helpers.MakeErrorResponse("ERROR", 400, err.Error(), links))
	}

	tx = handler.likeService.Delete(ID, eventID)

	if tx != nil {
		if reflect.TypeOf(tx).String() == "web.WebError" {
			webErr := tx.(web.WebError)
			return c.JSON(webErr.Code, helpers.MakeErrorResponse("ERROR", webErr.Code, webErr.Error(), links))
		}
	}

	return c.JSON(200, web.SuccessResponse{
		Status: "OK",
		Code:   200,
		Error:  nil,
		Links:  links,
		Data:   "Success dislike this event",
	})
}
