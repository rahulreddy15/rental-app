package service

import (
	"backend/internal/repository"

	"gorm.io/gorm"
)

type Services struct {
	User UserService
	db   *gorm.DB
}

func NewServices(db *gorm.DB, repos *repository.Repositories) *Services {
	return &Services{
		User: NewUserService(db, repos.User),
		db:   db,
	}
}

func (s *Services) Transaction(fn func(txServices *Services) error) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		txRepos := repository.NewRepositories(tx)
		txServices := NewServices(tx, txRepos)
		return fn(txServices)
	})
}
