package handler

import (
	"strconv"

	"backend/internal/model"
	"backend/internal/service"
	"backend/pkg/response"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type ListUsersResponse struct {
	Users  []model.User `json:"users"`
	Total  int64        `json:"total"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
}

// ListUsers godoc
// @Summary List all users
// @Description Get a paginated list of all users
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response{data=ListUsersResponse}
// @Router /users [get]
func (h *UserHandler) ListUsers(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}

	users, total, err := h.userService.List(c.Request().Context(), limit, offset)
	if err != nil {
		return response.FromError(c, err)
	}

	return response.Success(c, ListUsersResponse{
		Users:  users,
		Total:  total,
		Limit:  limit,
		Offset: offset,
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
// @Failure 409 {object} response.ErrorResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(c echo.Context) error {
	req := new(model.CreateUserRequest)

	if err := c.Bind(req); err != nil {
		return response.BadRequest(c, "Invalid request body", nil)
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	user, err := h.userService.Create(c.Request().Context(), service.CreateUserInput{
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	})
	if err != nil {
		return response.FromError(c, err)
	}

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
func (h *UserHandler) GetUser(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid user ID format", nil)
	}

	user, err := h.userService.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.FromError(c, err)
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
func (h *UserHandler) UpdateUser(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid user ID format", nil)
	}

	req := new(model.UpdateUserRequest)
	if err := c.Bind(req); err != nil {
		return response.BadRequest(c, "Invalid request body", nil)
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	input := service.UpdateUserInput{}
	if req.Name != "" {
		input.Name = &req.Name
	}
	if req.Role != "" {
		input.Role = &req.Role
	}

	user, err := h.userService.Update(c.Request().Context(), id, input)
	if err != nil {
		return response.FromError(c, err)
	}

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
func (h *UserHandler) DeleteUser(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid user ID format", nil)
	}

	if err := h.userService.Delete(c.Request().Context(), id); err != nil {
		return response.FromError(c, err)
	}

	return response.NoContent(c)
}
