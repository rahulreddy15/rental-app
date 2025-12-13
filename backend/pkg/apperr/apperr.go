package apperr

import "fmt"

type Code string

const (
	CodeNotFound     Code = "not_found"
	CodeConflict     Code = "conflict"
	CodeInvalid      Code = "invalid"
	CodeInternal     Code = "internal"
	CodeUnauthorized Code = "unauthorized"
	CodeForbidden    Code = "forbidden"
)

type AppError struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code Code, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

func NotFound(message string, err error) *AppError {
	return New(CodeNotFound, message, err)
}

func Conflict(message string, err error) *AppError {
	return New(CodeConflict, message, err)
}

func Invalid(message string, err error) *AppError {
	return New(CodeInvalid, message, err)
}

func Internal(message string, err error) *AppError {
	return New(CodeInternal, message, err)
}

func Unauthorized(message string, err error) *AppError {
	return New(CodeUnauthorized, message, err)
}

func Forbidden(message string, err error) *AppError {
	return New(CodeForbidden, message, err)
}
