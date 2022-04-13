package main

import (
	"tupulung/config"
	"tupulung/deliveries/handlers"
	"tupulung/deliveries/routes"
	"tupulung/utilities"

	categoryRepository "tupulung/repositories/category"
	commentRepository "tupulung/repositories/comment"
	eventRepository "tupulung/repositories/event"
	userRepository "tupulung/repositories/user"
	authService "tupulung/services/auth"
	categoryService "tupulung/services/category"
	commentService "tupulung/services/comment"
	eventService "tupulung/services/event"
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

	// Comment
	commentRepository := commentRepository.NewCommentRepository(db)
	commentService := commentService.NewCommentService(commentRepository)
	commentHandler := handlers.NewCommentHandler(commentService)
	routes.RegisterCommentRoute(e, commentHandler)

	e.Logger.Fatal(e.Start(":" + config.App.Port))
}
