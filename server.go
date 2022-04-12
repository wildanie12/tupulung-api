package main

import (
	"tupulung/config"
	"tupulung/deliveries/handlers"
	"tupulung/deliveries/routes"
	"tupulung/utilities"

	userRepository "tupulung/repositories/user"
	userService "tupulung/services/user"

	"github.com/labstack/echo/v4"
)

func main() {
	config := config.Get()
	db := utilities.NewMysqlGorm(config)
	utilities.Migrate(db)

	e := echo.New()

	// User 
	userRepository := userRepository.NewUserRepository(db)
	userService := userService.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)
	routes.RegisterUserRoute(e, userHandler)

	// Authentication
	authHandler := handlers.NewAuthHandler()
	routes.RegisterAuthRoute(e, authHandler)

	e.Logger.Fatal(e.Start(":" + config.App.Port))
}