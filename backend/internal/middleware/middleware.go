package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Setup configures all middleware for the Echo instance
func Setup(e *echo.Echo) {
	// Request ID for tracing
	e.Use(middleware.RequestID())

	// Logger
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} | ${status} | ${latency_human} | ${remote_ip} | ${method} ${uri}\n",
	}))

	// Recover from panics
	e.Use(middleware.Recover())

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Security headers (skip for Swagger)
	e.Use(secureWithSwaggerSkip())
}

// secureWithSwaggerSkip applies security headers but skips them for Swagger routes
func secureWithSwaggerSkip() echo.MiddlewareFunc {
	secureMiddleware := middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
	})

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip security headers for Swagger
			if strings.HasPrefix(c.Request().URL.Path, "/swagger") {
				return next(c)
			}
			return secureMiddleware(next)(c)
		}
	}
}
