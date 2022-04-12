package routes

import (
	"tupulung/deliveries/handlers"
	"tupulung/deliveries/middleware"

	"github.com/labstack/echo/v4"
)


func RegisterUserRoute(e *echo.Echo, userHandler *handlers.UserHandler) {
	group := e.Group("/api/users")
	group.POST("", userHandler.Create) 			// Registration
	group.GET("/:id", userHandler.Show)			// Detail User
	group.PUT("/:id", userHandler.Update, middleware.JWTMiddleware())		// Edit profile user
	group.DELETE("/:id", userHandler.Delete, middleware.JWTMiddleware())	// Delete account
}

func RegisterAuthRoute(e *echo.Echo, authHandler *handlers.AuthHandler)  {
	e.GET("/", authHandler.Index)
}