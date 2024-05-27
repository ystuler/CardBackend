package service

import (
	"back/internal/models"
	"back/internal/repository"
	"back/internal/schemas"
)

type Authorization interface {
	CreateUser(userSchema *schemas.CreateUserReq) (*models.User, error)
}

type Service struct {
	Authorization
}

func NewService(repo repository.UserRepository) *Service {
	return &Service{
		Authorization: NewUserService(repo),
	}
}
