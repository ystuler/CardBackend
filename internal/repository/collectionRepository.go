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

func (r *CollectionRepositoryImpl) GetCollectionByID(id int) (*models.Collection, error) {
	var collection models.Collection
	if err := r.db.First(&collection, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &collection, nil
}

func (r *CollectionRepositoryImpl) UpdateCollection(collection *models.Collection) (*models.Collection, error) {
	result := r.db.Model(collection).Updates(collection)
	if result.Error != nil {
		return nil, result.Error
	}
	return collection, nil
}

func (r *CollectionRepositoryImpl) RemoveCollection(collection *models.Collection) error {
	if err := r.db.Delete(&collection).Error; err != nil {
		return err
	}
	return nil
}

func (r *CollectionRepositoryImpl) GetAllCollections(userID int) (*[]models.Collection, error) {
	var collections []models.Collection
	if err := r.db.Where("user_id = ?", userID).Find(&collections).Error; err != nil {
		return nil, err
	}
	return &collections, nil
}
