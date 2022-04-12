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
	group.PUT("/:id", userHandler.Update, middleware.JWTMiddleware())    // Edit profile user
	group.DELETE("/:id", userHandler.Delete, middleware.JWTMiddleware()) // Delete account
}
func RegisterEventRoute(e *echo.Echo, eventHandler *handlers.EventHandler) {
	group := e.Group("/api/events")
	group.POST("", eventHandler.Create, middleware.JWTMiddleware())       // Registration event
	group.GET("/:id", eventHandler.Show)                                  // Detail event
	group.GET("/", eventHandler.Index)                                    // Get all Event
	e.GET("/api/users/:id/events", eventHandler.GetUserEvent)             // Detail event user
	group.PUT("/:id", eventHandler.Update, middleware.JWTMiddleware())    // Edit profile event
	group.DELETE("/:id", eventHandler.Delete, middleware.JWTMiddleware()) // Delete event
}

func RegisterAuthRoute(e *echo.Echo, authHandler *handlers.AuthHandler) {
	e.POST("/api/auth", authHandler.Login)
	e.GET("/api/auth/me", authHandler.Me, middleware.JWTMiddleware())
}
