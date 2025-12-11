package handler

import (
	"backend/pkg/response"

	"github.com/labstack/echo/v4"
)

// HealthResponse represents health check response
type HealthResponse struct {
	Status string `json:"status"`
}

// HealthCheck godoc
// @Summary Health check
// @Description Check if the API is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=HealthResponse}
// @Router /health [get]
func HealthCheck(c echo.Context) error {
	return response.Success(c, HealthResponse{
		Status: "healthy",
	})
}
