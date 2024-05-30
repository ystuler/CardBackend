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

func (s *CardServiceImpl) UpdateCard(cardSchema *schemas.UpdateCardReq) (*schemas.UpdateCardResp, error) {
	card, err := s.repo.GetCardByID(cardSchema.ID)
	if err != nil {
		return nil, err
	}

	if cardSchema.Question != nil {
		card.Question = *cardSchema.Question
	}

	if cardSchema.Answer != nil {
		card.Answer = *cardSchema.Answer
	}

	newCard, err := s.repo.UpdateCard(card)
	if err != nil {
		return nil, err
	}
	updatedCard := schemas.UpdateCardResp{
		ID:           newCard.ID,
		Question:     newCard.Question,
		Answer:       newCard.Answer,
		CollectionID: newCard.CollectionID,
	}
	return &updatedCard, nil
}

func (s *CardServiceImpl) RemoveCard(cardSchema *schemas.RemoveCardReq) error {
	card, err := s.repo.GetCardByID(cardSchema.ID)
	if err != nil {
		return err
	}
	err = s.repo.RemoveCard(card)
	if err != nil {
		return err
	}
	return nil
}
