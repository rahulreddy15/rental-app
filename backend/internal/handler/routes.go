package handler

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all API routes
func RegisterRoutes(g *echo.Group) {
	// Health check
	g.GET("/health", HealthCheck)

	// User routes
	users := g.Group("/users")
	{
		users.GET("", ListUsers)
		users.POST("", CreateUser)
		users.GET("/:id", GetUser)
		users.PUT("/:id", UpdateUser)
		users.DELETE("/:id", DeleteUser)
	}
}
