package service

import (
	"back/internal/repository"
	"back/internal/schemas"
)

type Authorization interface {
	CreateUser(userSchema *schemas.CreateUserReq) (*schemas.CreateUserResp, error)
	SignIn(userSchema *schemas.SignInReq) (*schemas.SignInResp, error)
}

type Service struct {
	Authorization
}

func NewService(repo repository.UserRepository) *Service {
	return &Service{
		Authorization: NewUserService(repo),
	}
}
