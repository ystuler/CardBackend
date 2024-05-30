package service

import (
	"back/internal/repository"
	"back/internal/schemas"
)

type Authorization interface {
	SignUp(userSchema *schemas.CreateUserReq) (*schemas.CreateUserResp, error)
	SignIn(userSchema *schemas.SignInReq) (*schemas.SignInResp, error)
}

type Collection interface {
	CreateCollection(collectionSchema *schemas.CreateCollectionReq, userID int) (*schemas.CreateCollectionResp, error)
	UpdateCollection(collectionSchema *schemas.UpdateCollectionReq, userID int) (*schemas.UpdateCollectionResp, error)
	RemoveCollection(collectionSchema *schemas.RemoveCollectionReq) error
	GetAllCollections(userID int) (*schemas.AllCollectionsResp, error)
}

type Service struct {
	Authorization
	Collection
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.UserRepository),
		Collection:    NewCollectionService(repos.CollectionRepository),
	}
}
