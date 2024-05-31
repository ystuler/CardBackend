package repository

import (
	"back/internal/models"
	"gorm.io/gorm"
)

type CardRepositoryImpl struct {
	db *gorm.DB
}

func NewCardRepo(db *gorm.DB) *CardRepositoryImpl {
	return &CardRepositoryImpl{db: db}
}

func (r *CardRepositoryImpl) CreateCard(card *models.Card) (*models.Card, error) {
	if err := r.db.Create(&card).Error; err != nil {
		return nil, err
	}
	return card, nil
}

func (r *CardRepositoryImpl) UpdateCard(card *models.Card) (*models.Card, error) {
	result := r.db.Model(card).Updates(card)
	if result.Error != nil {
		return nil, result.Error
	}
	return card, nil
}

func (r *CardRepositoryImpl) RemoveCard(card *models.Card) error {
	if err := r.db.Delete(&card).Error; err != nil {
		return err
	}
	return nil
}

func (r *CardRepositoryImpl) GetAllCards() ([]models.Card, error) {
	var cards []models.Card
	if err := r.db.Find(&cards).Error; err != nil {
		return nil, err
	}
	return cards, nil
}

func (r *CardRepositoryImpl) GetCardByID(cardID int) (*models.Card, error) {
	var card models.Card
	if err := r.db.First(&card, cardID).Error; err != nil {
		return nil, err
	}
	return &card, nil
}

func (r *CardRepositoryImpl) GetCardsByCollectionID(collectionID int) (*[]models.Card, error) {
	var cards []models.Card
	if err := r.db.Where("collection_id = ?", collectionID).Find(&cards).Error; err != nil {
		return nil, err
	}
	return &cards, nil
}
