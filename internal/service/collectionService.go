package service

import (
	"back/internal/exceptions"
	"back/internal/models"
	"back/internal/repository"
	"back/internal/schemas"
	"math/rand"
	"net/http"
	"time"
)

type CollectionServiceImpl struct {
	repo repository.CollectionRepository
}

func NewCollectionService(repo repository.CollectionRepository) *CollectionServiceImpl {
	return &CollectionServiceImpl{repo: repo}
}

func (s *CollectionServiceImpl) EnsureCollectionAccess(collectionID, userID int) error {
	collection, err := s.repo.GetCollectionByID(collectionID)
	if err != nil {
		return err
	}

	if collection.UserID != userID {
		return exceptions.NewAppError(http.StatusForbidden, exceptions.ErrForbidden, nil)
	}

	return nil
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

func (s *CollectionServiceImpl) GetCollectionByID(collectionID int) (*schemas.GetCollectionByIDResp, error) {
	collection, err := s.repo.GetCollectionByID(collectionID)
	if err != nil {
		return nil, err
	}

	description := ""
	if collection.Description != nil {
		description = *collection.Description
	}

	return &schemas.GetCollectionByIDResp{
		ID:          collection.ID,
		Name:        collection.Name,
		Description: description,
		CreatedAt:   collection.CreatedAt,
	}, nil
}

func (s *CollectionServiceImpl) UpdateCollection(collectionSchema *schemas.UpdateCollectionReq) (*schemas.UpdateCollectionResp, error) {
	collection, err := s.repo.GetCollectionByID(collectionSchema.ID)
	if err != nil {
		return nil, err
	}

	collection.Name = collectionSchema.Name
	collection.Description = &collectionSchema.Description

	newCollection, err := s.repo.UpdateCollection(collection)
	if err != nil {
		return nil, err
	}

	description := ""
	if newCollection.Description != nil {
		description = *newCollection.Description
	}

	updatedCollection := schemas.UpdateCollectionResp{
		ID:          newCollection.ID,
		Name:        newCollection.Name,
		Description: description,
		CreatedAt:   newCollection.CreatedAt,
	}
	return &updatedCollection, nil
}

func (s *CollectionServiceImpl) PatchCollection(collectionSchema *schemas.PatchCollectionReq) (*schemas.UpdateCollectionResp, error) {
	collection, err := s.repo.GetCollectionByID(collectionSchema.ID)
	if err != nil {
		return nil, err
	}

	if collectionSchema.Name != nil {
		collection.Name = *collectionSchema.Name
	}

	if collectionSchema.Description != nil {
		collection.Description = collectionSchema.Description
	}

	newCollection, err := s.repo.UpdateCollection(collection)
	if err != nil {
		return nil, err
	}

	description := ""
	if newCollection.Description != nil {
		description = *newCollection.Description
	}

	updatedCollection := schemas.UpdateCollectionResp{
		ID:          newCollection.ID,
		Name:        newCollection.Name,
		Description: description,
		CreatedAt:   newCollection.CreatedAt,
	}
	return &updatedCollection, nil
}

func (s *CollectionServiceImpl) RemoveCollection(collectionSchema *schemas.RemoveCollectionReq) error {
	collection, err := s.repo.GetCollectionByID(collectionSchema.ID)
	if err != nil {
		return err
	}

	err = s.repo.RemoveCollection(collection)
	if err != nil {
		return err
	}
	return nil
}

func (s *CollectionServiceImpl) GetAllCollections(userID int) (*schemas.AllCollectionsResp, error) {
	allCollections, err := s.repo.GetAllCollections(userID)
	if err != nil {
		return nil, err
	}

	collections := make([]schemas.AllCollections, len(*allCollections))
	for i, collection := range *allCollections {
		description := ""
		if collection.Description != nil {
			description = *collection.Description
		}

		collections[i] = schemas.AllCollections{
			ID:          collection.ID,
			Name:        collection.Name,
			Description: description,
			CreatedAt:   collection.CreatedAt,
		}
	}

	resp := &schemas.AllCollectionsResp{
		Collections: collections,
	}

	return resp, nil

}

func (s *CollectionServiceImpl) TrainCards(req *schemas.TrainSchemaReq) (*schemas.TrainSchemaResp, error) {
	cards, err := s.repo.GetAllCardsByCollectionID(req.ID)
	if err != nil {
		return nil, err
	}

	cardsSlice := *cards

	rand.Shuffle(len(cardsSlice), func(i, j int) {
		cardsSlice[i], cardsSlice[j] = cardsSlice[j], cardsSlice[i]
	})

	randCardSchema := make([]schemas.CardsByCollectionID, len(cardsSlice))
	for i, card := range cardsSlice {
		randCardSchema[i] = schemas.CardsByCollectionID{
			ID:       card.ID,
			Question: card.Question,
			Answer:   card.Answer,
		}
	}

	resp := schemas.TrainSchemaResp{Cards: randCardSchema}

	return &resp, nil
}
