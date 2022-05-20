package main

import (
	"tupulung/config"
	"tupulung/deliveries/handlers"
	"tupulung/deliveries/routes"
	"tupulung/utilities"

	categoryRepository "tupulung/repositories/category"
	commentRepository "tupulung/repositories/comment"
	eventRepository "tupulung/repositories/event"
	likeRepository "tupulung/repositories/like"
	participantRepository "tupulung/repositories/participant"
	userRepository "tupulung/repositories/user"
	authService "tupulung/services/auth"
	categoryService "tupulung/services/category"
	commentService "tupulung/services/comment"
	eventService "tupulung/services/event"
	likeService "tupulung/services/like"
	participantService "tupulung/services/participant"
	userService "tupulung/services/user"
	storageProvider "tupulung/utilities/storage"

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
	s3 := storageProvider.NewS3()

	// User
	userRepository := userRepository.NewUserRepository(db)

	eventRepository := eventRepository.NewEventRepository(db)
	likeRepository := likeRepository.NewLikeRepository(db)
	participantRepository := participantRepository.NewParticipantRepository(db)

	userService := userService.NewUserService(userRepository, eventRepository)
	userHandler := handlers.NewUserHandler(userService, s3)
	routes.RegisterUserRoute(e, userHandler)

	eventService := eventService.NewEventService(eventRepository, userRepository, likeRepository)
	participantService := participantService.NewParticipantService(participantRepository, userRepository, eventRepository)
	likeService := likeService.NewLikeService(likeRepository, userRepository, eventRepository)

	eventHandler := handlers.NewEventHandler(eventService, s3)
	participantHandler := handlers.NewParticipantHandler(participantService)
	likeHandler := handlers.NewLikeHandler(likeService)
	routes.RegisterEventRoute(e, eventHandler, participantHandler, likeHandler)

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
