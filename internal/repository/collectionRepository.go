package repository

import (
	"back/internal/models"
	"gorm.io/gorm"
)

type CollectionRepositoryImpl struct {
	db *gorm.DB
}

func NewCollectionRepo(db *gorm.DB) *CollectionRepositoryImpl {
	return &CollectionRepositoryImpl{db: db}
}

func (r *CollectionRepositoryImpl) CreateCollection(collection *models.Collection) (*models.Collection, error) {
	if err := r.db.Create(&collection).Error; err != nil {
		return nil, err
	}
	return collection, nil
}
