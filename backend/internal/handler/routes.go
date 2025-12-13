package handler

import (
	"backend/internal/service"

	"github.com/labstack/echo/v4"
)

type Handlers struct {
	User *UserHandler
}

func NewHandlers(services *service.Services) *Handlers {
	return &Handlers{
		User: NewUserHandler(services.User),
	}
}

func RegisterRoutes(g *echo.Group, handlers *Handlers) {
	g.GET("/health", HealthCheck)

	users := g.Group("/users")
	{
		users.GET("", handlers.User.ListUsers)
		users.POST("", handlers.User.CreateUser)
		users.GET("/:id", handlers.User.GetUser)
		users.PUT("/:id", handlers.User.UpdateUser)
		users.DELETE("/:id", handlers.User.DeleteUser)
	}
}
