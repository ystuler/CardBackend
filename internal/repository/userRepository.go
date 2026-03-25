package repository

import (
	"back/internal/exceptions"
	"back/internal/models"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) CreateUser(user *models.User) (*models.User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepositoryImpl) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username= ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewAppError(http.StatusNotFound, "user not found", err)
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetUserById(userId int) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewAppError(http.StatusNotFound, "user not found", err)
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) UpdateUser(user *models.User) (*models.User, error) {
	result := r.db.Model(user).Updates(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
