package config

import (
	"p2h-api/app/services"
)

type ServiceRegistry struct {
	UserService *services.UserService
}

func NewServiceRegistry(repository *RepositoryRegistry) *ServiceRegistry {
	return &ServiceRegistry{
		UserService: services.NewUserService(repository.UserRepo),
	}
}
