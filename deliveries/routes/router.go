package routes

import (
	"tupulung/deliveries/handlers"
	"tupulung/deliveries/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterUserRoute(e *echo.Echo, userHandler *handlers.UserHandler) {
	group := e.Group("/api/users")
	group.POST("", userHandler.Create)                                   // Registration
	group.GET("/:id", userHandler.Show)                                  // Detail User
	group.GET("/events", userHandler.GetUserEvents, middleware.JWTMiddleware())                                  // Detail User
	group.PUT("/:id", userHandler.Update, middleware.JWTMiddleware())    // Edit profile user
	group.DELETE("/:id", userHandler.Delete, middleware.JWTMiddleware()) // Delete account
}
func RegisterEventRoute(e *echo.Echo, eventHandler *handlers.EventHandler, participantHandler *handlers.ParticipantHandler, likeHandler *handlers.LikeHandler) {
	group := e.Group("/api/events")
	group.POST("", eventHandler.Create, middleware.JWTMiddleware())                   // Registration event
	group.GET("", eventHandler.Index)                                                 // Get all Event
	group.GET("/:id", eventHandler.Show)                                              // Detail event
	e.GET("/api/users/:id/events", eventHandler.GetUserEvent)                         // Detail user event
	group.PUT("/:id", eventHandler.Update, middleware.JWTMiddleware())                // Edit profile event
	group.DELETE("/:id", eventHandler.Delete, middleware.JWTMiddleware())             // Delete event
	group.POST("/join/:id", participantHandler.Append, middleware.JWTMiddleware())    // Join an event
	group.DELETE("/leave/:id", participantHandler.Delete, middleware.JWTMiddleware()) // Leave an event
	group.POST("/like/:id", likeHandler.Append, middleware.JWTMiddleware())           // Like an event
	group.DELETE("/dislike/:id", likeHandler.Delete, middleware.JWTMiddleware())      // Dislike an event
}

func RegisterAuthRoute(e *echo.Echo, authHandler *handlers.AuthHandler) {
	e.POST("/api/auth", authHandler.Login)
	e.GET("/api/auth/me", authHandler.Me, middleware.JWTMiddleware())
}

func RegisterCategoryRoute(e *echo.Echo, categoryHandler handlers.CategoryHandler) {
	e.GET("/api/categories", categoryHandler.Index)
	e.POST("/api/categories", categoryHandler.Create, middleware.JWTMiddleware())
	e.PUT("/api/categories/:id", categoryHandler.Update, middleware.JWTMiddleware())
	e.DELETE("/api/categories/:id", categoryHandler.Delete, middleware.JWTMiddleware())
}

func RegisterCommentRoute(e *echo.Echo, commentHandler *handlers.CommentHandler) {
	e.GET("/api/events/:eventID/comments", commentHandler.Index)
	e.POST("/api/events/:eventID/comments", commentHandler.Create, middleware.JWTMiddleware())
	e.PUT("/api/events/comments/:commentID", commentHandler.Update, middleware.JWTMiddleware())
	e.DELETE("/api/events/comments/:commentID", commentHandler.Delete, middleware.JWTMiddleware())
}
