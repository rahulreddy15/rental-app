package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}

// Success returns a success response
func Success(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// Created returns a 201 created response
func Created(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

// NoContent returns a 204 no content response
func NoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

// Error returns an error response
func Error(c echo.Context, status int, message string, details interface{}) error {
	return c.JSON(status, ErrorResponse{
		Success: false,
		Error:   message,
		Details: details,
	})
}

// BadRequest returns a 400 bad request response
func BadRequest(c echo.Context, message string, details interface{}) error {
	return Error(c, http.StatusBadRequest, message, details)
}

// NotFound returns a 404 not found response
func NotFound(c echo.Context, message string) error {
	return Error(c, http.StatusNotFound, message, nil)
}

// InternalError returns a 500 internal server error response
func InternalError(c echo.Context, message string) error {
	return Error(c, http.StatusInternalServerError, message, nil)
}
