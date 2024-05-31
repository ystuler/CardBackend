package schemas

import "time"

//todo у id может быть только положительные id, поэтоум заменить на uint

type CreateCollectionReq struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
}

type CreateCollectionResp struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

type UpdateCollectionReq struct {
	ID          int    `validate:"required,gt=0"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type UpdateCollectionResp struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

type RemoveCollectionReq struct {
	ID int `validate:"required,gt=0"`
}

type AllCollectionsReq struct {
	ID int `validate:"required,gt=0"`
}
type AllCollections struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type AllCollectionsResp struct {
	Collections []AllCollections `json:"collections"`
}

type TrainSchemaReq struct {
	ID int `validate:"required,gt=0"`
}

type TrainSchemaResp struct {
	Cards []CardsByCollectionID `json:"cards"`
}
