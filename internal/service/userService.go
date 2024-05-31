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

type AuthenticationImpl struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *AuthenticationImpl {
	return &AuthenticationImpl{repo: repo}
}

func (s *AuthenticationImpl) SignUp(userSchema *schemas.CreateUserReq) (*schemas.CreateUserResp, error) {
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

	generatedJWT, err := util.GenerateJWT(createdUser)
	if err != nil {
		return nil, err
	}

	userSchemaResp := schemas.CreateUserResp{
		Token:    generatedJWT,
		ID:       createdUser.ID,
		Username: createdUser.Username,
	}

	return &userSchemaResp, nil
}

func (s *AuthenticationImpl) SignIn(userSchema *schemas.SignInReq) (*schemas.SignInResp, error) {
	existingUser, err := s.repo.GetUserByUsername(userSchema.Username)

	//fixme не срабатывает ошибка
	//if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
	//	return nil, errors.New("user not found")
	//}

	if err != nil {
		return nil, err
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

func (s *AuthenticationImpl) GetProfile(userID int) (*schemas.GetProfileResp, error) {
	user, err := s.repo.GetUserById(userID)
	if err != nil {
		return nil, err
	}

	profile := schemas.Profile{
		ID:       user.ID,
		Username: user.Username,
	}

	return &schemas.GetProfileResp{Profile: profile}, nil
}

func (s *AuthenticationImpl) UpdateUsername(userID int, usernameSchema *schemas.UpdateUsernameReq) (*schemas.UpdateUsernameResp, error) {
	user, err := s.repo.GetUserById(userID)
	if err != nil {
		return nil, err
	}

	user.Username = usernameSchema.Username

	updatedUser, err := s.repo.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return &schemas.UpdateUsernameResp{ID: updatedUser.ID, Username: updatedUser.Username}, nil
}
