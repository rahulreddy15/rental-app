package validator

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator wraps go-playground/validator for Echo
type CustomValidator struct {
	validator *validator.Validate
}

// NewValidator creates a new CustomValidator instance
func NewValidator() *CustomValidator {
	v := validator.New()

	// Use JSON tag names in error messages instead of struct field names
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom validations here
	// Example: v.RegisterValidation("custom_tag", customValidationFunc)

	return &CustomValidator{validator: v}
}

// Validate implements echo.Validator interface
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, formatValidationErrors(err))
	}
	return nil
}

// ValidationErrorResponse represents validation error details
type ValidationErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// formatValidationErrors converts validator errors to readable format
func formatValidationErrors(err error) []ValidationErrorResponse {
	var errors []ValidationErrorResponse

	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, ValidationErrorResponse{
			Field:   e.Field(),
			Message: getErrorMessage(e),
		})
	}

	return errors
}

// getErrorMessage returns human-readable error message for validation tag
func getErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short, minimum is " + e.Param()
	case "max":
		return "Value is too long, maximum is " + e.Param()
	case "gte":
		return "Value must be greater than or equal to " + e.Param()
	case "lte":
		return "Value must be less than or equal to " + e.Param()
	case "oneof":
		return "Value must be one of: " + e.Param()
	default:
		return "Validation failed on " + e.Tag()
	}
}
