package handlers

import (
	"net/http"
	"reflect"
	"strconv"
	"tupulung/config"
	"tupulung/deliveries/helpers"
	"tupulung/deliveries/middleware"
	"tupulung/entities/web"
	participantService "tupulung/services/participant"

	"github.com/labstack/echo/v4"
)

type ParticipantHandler struct {
	participantService *participantService.ParticipantService
}

func NewParticipantHandler(service *participantService.ParticipantService) *ParticipantHandler {
	return &ParticipantHandler{
		participantService: service,
	}
}

func (handler ParticipantHandler) Append(c echo.Context) error {

	eventID, err := strconv.Atoi(c.Param("id"))
	links := map[string]string{"self": config.Get().App.BaseURL + "/events/" + c.Param("id")}
	token := c.Get("user")
	userID, err := middleware.ReadToken(token)
	if err != nil {
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

	tx := handler.participantService.Append(userID, eventID)

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
		Data:   "Success joined this event",
	})
}

func (handler ParticipantHandler) Delete(c echo.Context) error {

	eventID, err := strconv.Atoi(c.Param("id"))
	links := map[string]string{"self": config.Get().App.BaseURL + "/events/" + c.Param("id")}
	token := c.Get("user")
	userID, err := middleware.ReadToken(token)
	if err != nil {
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

	tx := handler.participantService.Delete(userID, eventID)

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
		Data:   "Success leave this event",
	})
}
