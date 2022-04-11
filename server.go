package main

import (
	"tupulung/config"
	"tupulung/deliveries/handlers"
	"tupulung/deliveries/routes"
	"tupulung/utilities"

	"github.com/labstack/echo/v4"
)

func main() {
	config := config.Get()
	utilities.NewMysqlGorm(config)

	e := echo.New()

	authHandler := handlers.NewAuthHandler()
	routes.RegisterAuthRoute(e, authHandler)

	e.Logger.Fatal(e.Start(":" + config.App.Port))
}