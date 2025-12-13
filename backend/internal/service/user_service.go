package service

import (
	"context"
	"errors"
	"time"

	"backend/internal/model"
	"backend/internal/repository"
	"backend/pkg/apperr"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	List(ctx context.Context, limit, offset int) ([]model.User, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	Create(ctx context.Context, input CreateUserInput) (*model.User, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateUserInput) (*model.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateUserInput struct {
	Name  string
	Email string
	Role  string
}

type UpdateUserInput struct {
	Name *string
	Role *string
}

type userService struct {
	db       *gorm.DB
	userRepo repository.UserRepository
}

func NewUserService(db *gorm.DB, userRepo repository.UserRepository) UserService {
	return &userService{
		db:       db,
		userRepo: userRepo,
	}
}

func (s *userService) List(ctx context.Context, limit, offset int) ([]model.User, int64, error) {
	users, total, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, apperr.Internal("Failed to fetch users", err)
	}
	return users, total, nil
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, apperr.NotFound("User not found", err)
		}
		return nil, apperr.Internal("Failed to fetch user", err)
	}
	return user, nil
}

func (s *userService) Create(ctx context.Context, input CreateUserInput) (*model.User, error) {
	user := &model.User{
		ID:        uuid.New(),
		Name:      input.Name,
		Email:     input.Email,
		Role:      input.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			return nil, apperr.Conflict("User with this email already exists", err)
		}
		return nil, apperr.Internal("Failed to create user", err)
	}

	return user, nil
}

func (s *userService) Update(ctx context.Context, id uuid.UUID, input UpdateUserInput) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, apperr.NotFound("User not found", err)
		}
		return nil, apperr.Internal("Failed to fetch user", err)
	}

	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Role != nil {
		user.Role = *input.Role
	}
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, apperr.Internal("Failed to update user", err)
	}

	return user, nil
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.userRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return apperr.NotFound("User not found", err)
		}
		return apperr.Internal("Failed to delete user", err)
	}
	return nil
}
