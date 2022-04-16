package handlers

import (
	"reflect"
	"strconv"
	"tupulung/config"
	"tupulung/deliveries/helpers"
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

	eventID, err := strconv.Atoi(c.Param("id"))
	links := map[string]string{"self": config.Get().App.BaseURL + "/events/" + c.Param("id")}
	if err != nil {
		return c.JSON(400, helpers.MakeErrorResponse("ERROR", 400, err.Error(), links))
	}

	tx := handler.likeService.Append(token, eventID)

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

	eventID, err := strconv.Atoi(c.Param("id"))
	links := map[string]string{"self": config.Get().App.BaseURL + "/events/" + c.Param("id")}
	if err != nil {
		return c.JSON(400, helpers.MakeErrorResponse("ERROR", 400, err.Error(), links))
	}

	tx := handler.likeService.Delete(token, eventID)

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