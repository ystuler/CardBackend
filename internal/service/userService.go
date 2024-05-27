package service

import (
	"back/internal/models"
	"back/internal/repository"
	"back/internal/schemas"
	"back/internal/util"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type UserServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{repo: repo}
}

func (s *UserServiceImpl) CreateUser(userSchema *schemas.CreateUserReq) (*models.User, error) {
	existingUser, err := s.repo.GetUserByUsername(userSchema.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := util.HashPassword(userSchema.Password)

	//todo create secret word
	user := &models.User{
		Username:     userSchema.Username,
		PasswordHash: hashedPassword,
		SecretWord:   "",
		CreatedAt:    time.Time{},
	}

	createdUser, err := s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	fmt.Println(createdUser.ID)
	fmt.Println(createdUser.Username)
	return createdUser, nil
}
