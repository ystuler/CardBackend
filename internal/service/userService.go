package service

import (
	"back/internal/models"
	"back/internal/repository"
	"back/internal/schemas"
	"back/internal/util"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"time"
)

type UserServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{repo: repo}
}

const signingKey = "mySecret"

func (s *UserServiceImpl) generateJWT(userModel *models.User) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   userModel.Username,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(signingKey))
}

func (s *UserServiceImpl) parseToken(accessToken string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
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

	generatedJWT, err := s.generateJWT(existingUser)
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
