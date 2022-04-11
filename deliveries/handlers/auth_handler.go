package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {

}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (handler AuthHandler) Index(c echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}