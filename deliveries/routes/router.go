package routes

import (
	"tupulung/deliveries/handlers"

	"github.com/labstack/echo/v4"
)


func RegisterAuthRoute(e *echo.Echo, authHandler *handlers.AuthHandler)  {
	e.GET("/", authHandler.Index)
}