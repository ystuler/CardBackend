package service

import (
	"back/internal/models"
	"back/internal/repository"
	"back/internal/schemas"
	"time"
)

type CollectionServiceImpl struct {
	repo repository.CollectionRepository
}

func NewCollectionService(repo repository.CollectionRepository) *CollectionServiceImpl {
	return &CollectionServiceImpl{repo: repo}
}

func (s *CollectionServiceImpl) CreateCollection(collectionSchema *schemas.CreateCollectionReq, userID int) (*schemas.CreateCollectionResp, error) {
	collection := &models.Collection{
		Name:        collectionSchema.Name,
		Description: &collectionSchema.Description,
		CreatedAt:   time.Now(),
		UserID:      userID,
	}
	createdCollection, err := s.repo.CreateCollection(collection)
	if err != nil {
		return nil, err
	}
	collectionResp := schemas.CreateCollectionResp{
		ID:          createdCollection.ID,
		Name:        createdCollection.Name,
		Description: *createdCollection.Description,
		CreatedAt:   createdCollection.CreatedAt,
	}

	return &collectionResp, nil
}
