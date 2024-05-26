package service

import (
	"back/internal/models"
	"back/internal/repository"
)

type UserService interface {
	Create(user *models.User) (*models.User, error)
}

type Service struct {
	UserService
}

func NewService(repo repository.UserRepository) *Service {
	return &Service{
		UserService: NewUserService(repo),
	}
}
