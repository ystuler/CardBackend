package service

import (
	"back/internal/exceptions"
	"back/internal/models"
	"back/internal/repository"
	"back/internal/schemas"
	"net/http"
)

type CardServiceImpl struct {
	cardRepo       repository.CardRepository
	collectionRepo repository.CollectionRepository
}

func NewCardService(cardRepo repository.CardRepository, collectionRepo repository.CollectionRepository) *CardServiceImpl {
	return &CardServiceImpl{cardRepo: cardRepo, collectionRepo: collectionRepo}
}

func (s *CardServiceImpl) EnsureCardAccess(cardID, userID int) error {
	card, err := s.cardRepo.GetCardByID(cardID)
	if err != nil {
		return err
	}

	collection, err := s.collectionRepo.GetCollectionByID(card.CollectionID)
	if err != nil {
		return err
	}

	if collection.UserID != userID {
		return exceptions.NewAppError(http.StatusForbidden, exceptions.ErrForbidden, nil)
	}

	return nil
}

func (s *CardServiceImpl) CreateCard(cardSchema *schemas.CreateCardReq, collectionID int) (*schemas.CreateCardResp, error) {
	card := &models.Card{
		CollectionID: collectionID,
		Question:     cardSchema.Question,
		Answer:       cardSchema.Answer,
	}

	createdCard, err := s.cardRepo.CreateCard(card)
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
	card, err := s.cardRepo.GetCardByID(cardSchema.ID)
	if err != nil {
		return nil, err
	}

	if cardSchema.Question != nil {
		card.Question = *cardSchema.Question
	}

	if cardSchema.Answer != nil {
		card.Answer = *cardSchema.Answer
	}

	newCard, err := s.cardRepo.UpdateCard(card)
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
	card, err := s.cardRepo.GetCardByID(cardSchema.ID)
	if err != nil {
		return err
	}
	err = s.cardRepo.RemoveCard(card)
	if err != nil {
		return err
	}
	return nil
}

func (s *CardServiceImpl) GetCardsByCollectionID(collectionID int) (*schemas.GetCardByCollectionIDResp, error) {
	allCards, err := s.cardRepo.GetCardsByCollectionID(collectionID)
	if err != nil {
		return nil, err
	}

	cards := make([]schemas.CardsByCollectionID, len(*allCards))
	for i, card := range *allCards {
		cards[i] = schemas.CardsByCollectionID{
			ID:       card.ID,
			Question: card.Question,
			Answer:   card.Answer,
		}
	}
	resp := &schemas.GetCardByCollectionIDResp{
		CollectionID: collectionID,
		Cards:        cards,
	}
	return resp, nil
}
