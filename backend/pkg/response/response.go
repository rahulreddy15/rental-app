package response

import (
	"errors"
	"net/http"

	"backend/pkg/apperr"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Code    string      `json:"code,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

func Success(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

func Created(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

func NoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func Error(c echo.Context, status int, code, message string, details interface{}) error {
	return c.JSON(status, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    code,
		Details: details,
	})
}

func FromError(c echo.Context, err error) error {
	var ae *apperr.AppError
	if errors.As(err, &ae) {
		status := codeToStatus(ae.Code)
		c.Logger().Error(err)
		return c.JSON(status, ErrorResponse{
			Success: false,
			Error:   ae.Message,
			Code:    string(ae.Code),
		})
	}

	c.Logger().Error(err)
	return Error(c, http.StatusInternalServerError, "internal", "Internal server error", nil)
}

func codeToStatus(code apperr.Code) int {
	switch code {
	case apperr.CodeNotFound:
		return http.StatusNotFound
	case apperr.CodeConflict:
		return http.StatusConflict
	case apperr.CodeInvalid:
		return http.StatusBadRequest
	case apperr.CodeUnauthorized:
		return http.StatusUnauthorized
	case apperr.CodeForbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

func BadRequest(c echo.Context, message string, details interface{}) error {
	return Error(c, http.StatusBadRequest, "invalid", message, details)
}

func NotFound(c echo.Context, message string) error {
	return Error(c, http.StatusNotFound, "not_found", message, nil)
}

func InternalError(c echo.Context, message string) error {
	return Error(c, http.StatusInternalServerError, "internal", message, nil)
}
