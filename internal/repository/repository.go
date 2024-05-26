package repository

import (
	"back/internal/models"
	"gorm.io/gorm"
)

type User interface {
	CreateUser(user *models.User) (*models.User, error)
}

type Repository struct {
	User
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User: NewUserRepo(db),
	}
}
