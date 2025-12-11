package handler

import (
	"net/http"
	"time"

	"backend/internal/model"
	"backend/pkg/response"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// In-memory store for demo purposes
// Replace with your database layer
var users = make(map[string]*model.User)

// ListUsersResponse represents paginated user list
type ListUsersResponse struct {
	Users []model.User `json:"users"`
	Total int          `json:"total"`
}

// ListUsers godoc
// @Summary List all users
// @Description Get a list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=ListUsersResponse}
// @Router /users [get]
func ListUsers(c echo.Context) error {
	userList := make([]model.User, 0, len(users))
	for _, u := range users {
		userList = append(userList, *u)
	}

	return response.Success(c, ListUsersResponse{
		Users: userList,
		Total: len(userList),
	})
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param user body model.CreateUserRequest true "User details"
// @Success 201 {object} response.Response{data=model.User}
// @Failure 400 {object} response.ErrorResponse
// @Router /users [post]
func CreateUser(c echo.Context) error {
	req := new(model.CreateUserRequest)

	// Bind request body
	if err := c.Bind(req); err != nil {
		return response.BadRequest(c, "Invalid request body", nil)
	}

	// Validate request
	if err := c.Validate(req); err != nil {
		return err // Validator already returns proper HTTP error
	}

	// Create user
	user := &model.User{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Email:     req.Email,
		Role:      req.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	users[user.ID] = user

	return response.Created(c, user)
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Get user details by user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.Response{data=model.User}
// @Failure 404 {object} response.ErrorResponse
// @Router /users/{id} [get]
func GetUser(c echo.Context) error {
	id := c.Param("id")

	user, exists := users[id]
	if !exists {
		return response.NotFound(c, "User not found")
	}

	return response.Success(c, user)
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update user details by user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body model.UpdateUserRequest true "User update details"
// @Success 200 {object} response.Response{data=model.User}
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /users/{id} [put]
func UpdateUser(c echo.Context) error {
	id := c.Param("id")

	user, exists := users[id]
	if !exists {
		return response.NotFound(c, "User not found")
	}

	req := new(model.UpdateUserRequest)

	if err := c.Bind(req); err != nil {
		return response.BadRequest(c, "Invalid request body", nil)
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	// Update fields if provided
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	user.UpdatedAt = time.Now()

	return response.Success(c, user)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Failure 404 {object} response.ErrorResponse
// @Router /users/{id} [delete]
func DeleteUser(c echo.Context) error {
	id := c.Param("id")

	if _, exists := users[id]; !exists {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			Success: false,
			Error:   "User not found",
		})
	}

	delete(users, id)

	return response.NoContent(c)
}
