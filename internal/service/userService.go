package service

import (
	"back/internal/exceptions"
	"back/internal/models"
	"back/internal/repository"
	"back/internal/schemas"
	"back/internal/util"
	"errors"
	"net/http"
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
	if err != nil {
		var appErr *exceptions.AppError
		if !errors.As(err, &appErr) || appErr.StatusCode != http.StatusNotFound {
			return nil, err
		}
	}
	if existingUser != nil {
		return nil, exceptions.NewAppError(http.StatusUnprocessableEntity, exceptions.ErrUserAlreadyExists, nil)
	}

	hashedPassword, err := util.HashPassword(userSchema.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     userSchema.Username,
		PasswordHash: hashedPassword,
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

	if err != nil {
		var appErr *exceptions.AppError
		if errors.As(err, &appErr) && appErr.StatusCode == http.StatusNotFound {
			return nil, exceptions.NewAppError(http.StatusUnauthorized, exceptions.ErrInvalidCredentials, err)
		}
		return nil, err
	}

	err = util.CheckPassword(userSchema.Password, existingUser.PasswordHash)
	if err != nil {
		return nil, exceptions.NewAppError(http.StatusUnauthorized, exceptions.ErrInvalidCredentials, err)
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

func (s *AuthenticationImpl) UpdateUsername(usernameSchema *schemas.UpdateUsernameReq) (*schemas.UpdateUsernameResp, error) {
	user, err := s.repo.GetUserById(usernameSchema.ID)
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

func (s *AuthenticationImpl) UpdatePassword(passwordSchema *schemas.UpdatePasswordReq) error {
	user, err := s.repo.GetUserById(passwordSchema.ID)
	if err != nil {
		return err
	}

	if err := util.CheckPassword(passwordSchema.OldPassword, user.PasswordHash); err != nil {
		return exceptions.NewAppError(http.StatusUnprocessableEntity, "old password does not match", err)
	}

	hashedPassword, err := util.HashPassword(passwordSchema.NewPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword

	_, err = s.repo.UpdateUser(user)
	return err
}
