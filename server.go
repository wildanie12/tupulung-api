package main

import (
	"tupulung/config"
	"tupulung/deliveries/handlers"
	"tupulung/deliveries/routes"
	"tupulung/utilities"

	categoryRepository "tupulung/repositories/category"
	eventRepository "tupulung/repositories/event"
	userRepository "tupulung/repositories/user"
	authService "tupulung/services/auth"
	categoryService "tupulung/services/category"
	eventService "tupulung/services/event"
	userService "tupulung/services/user"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config := config.Get()
	db := utilities.NewMysqlGorm(config)
	utilities.Migrate(db)

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
	}))

	// User
	userRepository := userRepository.NewUserRepository(db)
	userService := userService.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)
	routes.RegisterUserRoute(e, userHandler)

	// Event
	eventRepository := eventRepository.NewEventRepository(db)
	eventService := eventService.NewEventService(eventRepository, userRepository)
	eventHandler := handlers.NewEventHandler(eventService)
	routes.RegisterEventRoute(e, eventHandler)

	// Authentication
	authService := authService.NewAuthService(userRepository)
	authHandler := handlers.NewAuthHandler(authService)
	routes.RegisterAuthRoute(e, authHandler)

	// User 
	categoryRepository := categoryRepository.NewCategoryRepository(db)
	categoryService := categoryService.NewCategoryService(categoryRepository)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	routes.RegisterCategoryRoute(e, categoryHandler)

	e.Logger.Fatal(e.Start(":" + config.App.Port))
}
