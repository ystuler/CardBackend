package service

import (
	"back/internal/models"
	"back/internal/repository"
	"back/internal/schemas"
	"math/rand"
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

// todo удаление коллекции (cascade)
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
		collections[i] = schemas.AllCollections{
			ID:          collection.ID,
			Name:        collection.Name,
			Description: *collection.Description,
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
