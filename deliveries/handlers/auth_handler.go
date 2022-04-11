package handlers

import (
	"net/http"

	"tupulung/config"
	"tupulung/entities/web"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {

}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (handler AuthHandler) Index(c echo.Context) error {
	// Define links
	links := map[string]string{}
	links["self"] = config.Get().App.BaseURL + "/"

	// Response
	return c.JSON(http.StatusOK, web.SuccessResponse {
		Status: "OK",
		Code: 400,
		Error: nil,
		Links: links,
	})
}