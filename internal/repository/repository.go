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

type CardRepository interface {
	CreateCard(card *models.Card) (*models.Card, error)
	UpdateCard(card *models.Card) (*models.Card, error)
	RemoveCard(card *models.Card) error
	GetAllCards() ([]models.Card, error)
	GetCardByID(cardID int) (*models.Card, error)
	GetCardsByCollectionID(collectionID int) ([]models.Card, error)
}

type Repository struct {
	UserRepository
	CollectionRepository
	CardRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		UserRepository:       NewUserRepo(db),
		CollectionRepository: NewCollectionRepo(db),
		CardRepository:       NewCardRepo(db),
	}
}
