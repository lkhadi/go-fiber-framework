package config

import (
	"p2h-api/app/repositories"

	"gorm.io/gorm"
)

type RepositoryRegistry struct {
	UserRepo repositories.UserRepositoryInterface
}

func NewRepositoryRegistry(db *gorm.DB) *RepositoryRegistry {
	return &RepositoryRegistry{
		UserRepo: repositories.NewUserRepository(db),
	}
}
