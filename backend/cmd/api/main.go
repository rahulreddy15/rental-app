// @title Rental Property Management API
// @version 1.0
// @description API for managing rental properties in India
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "backend/docs"
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/repository"
	"backend/internal/service"
	customValidator "backend/internal/validator"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.RunMigrations(&cfg.Database, "migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	repos := repository.NewRepositories(db)
	services := service.NewServices(db, repos)
	handlers := handler.NewHandlers(services)

	e := echo.New()
	e.Validator = customValidator.NewValidator()

	middleware.Setup(e)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	api := e.Group("/api/v1")
	handler.RegisterRoutes(api, handlers)

	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := e.Start(":" + cfg.Port); err != nil {
			log.Printf("Server stopped: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := database.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}
	log.Println("Server gracefully stopped")
}
