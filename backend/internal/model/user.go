package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name      string    `json:"name" gorm:"type:varchar(100);not null"`
	Email     string    `json:"email" gorm:"type:varchar(255);not null;uniqueIndex"`
	Role      string    `json:"role" gorm:"type:varchar(20);not null;default:'user'"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:now()"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (User) TableName() string {
	return "users"
}

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof=admin user guest"`
}

type UpdateUserRequest struct {
	Name string `json:"name" validate:"omitempty,min=2,max=100"`
	Role string `json:"role" validate:"omitempty,oneof=admin user guest"`
}
