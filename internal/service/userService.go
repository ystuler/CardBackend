package service

import (
	"back/internal/models"
	"back/internal/repository"
	"back/internal/schemas"
	"back/internal/util"
	"errors"
	"gorm.io/gorm"
	"time"
)

type UserServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{repo: repo}
}

func (s *UserServiceImpl) CreateUser(userSchema *schemas.CreateUserReq) (*schemas.CreateUserResp, error) {
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

	userSchemaResp := schemas.CreateUserResp{
		ID:       createdUser.ID,
		Username: createdUser.Username,
	}

	return &userSchemaResp, nil
}

func (s *UserServiceImpl) SignIn(userSchema *schemas.SignInReq) (*schemas.SignInResp, error) {
	existingUser, err := s.repo.GetUserByUsername(userSchema.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}

	err = util.CheckPassword(userSchema.Password, existingUser.PasswordHash)
	if err != nil {
		return nil, errors.New("password does not match")
	}

	generatedJWT, err := util.GenerateJWT(existingUser)
	if err != nil {
		return nil, err
	}

	userSchemaResp := schemas.SignInResp{
		Token:    generatedJWT,
		ID:       existingUser.ID,
		Username: existingUser.Username,
	}

	return &userSchemaResp, nil
}
