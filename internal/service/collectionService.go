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

func (s *CollectionServiceImpl) UpdateCollection(collectionSchema *schemas.UpdateCollectionReq, userID int) (*schemas.UpdateCollectionResp, error) {
	collection, err := s.repo.GetCollectionByID(collectionSchema.ID)
	if err != nil {
		return nil, err
	}

	//if collection == nil || collection.UserID != userID {
	//	return nil, errors.New("collection not found or unauthorized")
	//}

	collection.Name = collectionSchema.Name
	collection.Description = &collectionSchema.Description

	newCollection, err := s.repo.UpdateCollection(collection)
	if err != nil {
		return nil, err
	}

	updatedCollection := schemas.UpdateCollectionResp{
		ID:          newCollection.ID,
		Name:        newCollection.Name,
		Description: *newCollection.Description,
		CreatedAt:   newCollection.CreatedAt,
	}
	return &updatedCollection, nil
}
