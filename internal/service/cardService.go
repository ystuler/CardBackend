package service

import (
	"back/internal/models"
	"back/internal/repository"
	"back/internal/schemas"
)

type CardServiceImpl struct {
	repo repository.CardRepository
}

func NewCardService(repo repository.CardRepository) *CardServiceImpl {
	return &CardServiceImpl{repo: repo}
}

func (s *CardServiceImpl) CreateCard(cardSchema *schemas.CreateCardReq, collectionID int) (*schemas.CreateCardResp, error) {
	card := &models.Card{
		CollectionID: collectionID,
		Question:     cardSchema.Question,
		Answer:       cardSchema.Answer,
	}

	createdCard, err := s.repo.CreateCard(card)
	if err != nil {
		return nil, err
	}

	cardResp := &schemas.CreateCardResp{
		ID:           createdCard.ID,
		Question:     createdCard.Question,
		Answer:       createdCard.Answer,
		CollectionID: createdCard.CollectionID,
	}

	return cardResp, nil
}
