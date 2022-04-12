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