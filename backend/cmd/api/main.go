// @title Echo API Starter
// @version 1.0
// @description A production-ready Echo API starter template
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

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "backend/docs" // swagger docs
	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/middleware"
	customValidator "backend/internal/validator"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create Echo instance
	e := echo.New()

	// Set custom validator
	e.Validator = customValidator.NewValidator()

	// Middleware
	middleware.Setup(e)

	// Swagger endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// API routes
	api := e.Group("/api/v1")
	handler.RegisterRoutes(api)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
