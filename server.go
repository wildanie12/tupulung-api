package main

import (
	"tupulung/config"
	"tupulung/deliveries/handlers"
	"tupulung/deliveries/routes"
	"tupulung/utilities"

	categoryRepository "tupulung/repositories/category"
	commentRepository "tupulung/repositories/comment"
	eventRepository "tupulung/repositories/event"
	participantRepository "tupulung/repositories/participant"
	userRepository "tupulung/repositories/user"
	authService "tupulung/services/auth"
	categoryService "tupulung/services/category"
	commentService "tupulung/services/comment"
	eventService "tupulung/services/event"
	participantService "tupulung/services/participant"
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
	participantRepository := participantRepository.NewParticipantRepository(db)
	participantService := participantService.NewParticipantService(participantRepository, userRepository, eventRepository)
	participantHandler := handlers.NewParticipantHandler(participantService)
	routes.RegisterEventRoute(e, eventHandler, participantHandler)

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
	commentService := commentService.NewCommentService(commentRepository, userRepository)
	commentHandler := handlers.NewCommentHandler(commentService)
	routes.RegisterCommentRoute(e, commentHandler)

	// routes.RegisterParticipantRoute(e, participantHandler)

	e.Logger.Fatal(e.Start(":" + config.App.Port))
}
