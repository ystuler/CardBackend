package repository

import (
	"back/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
}

type CollectionRepository interface {
	CreateCollection(collection *models.Collection) (*models.Collection, error)
	GetCollectionByID(id int) (*models.Collection, error)
	UpdateCollection(collection *models.Collection) (*models.Collection, error)
}

type Repository struct {
	UserRepository
	CollectionRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		UserRepository:       NewUserRepo(db),
		CollectionRepository: NewCollectionRepo(db),
	}
}
