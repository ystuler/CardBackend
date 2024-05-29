package schemas

import "time"

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
	ID          int    `json:"id" validate:"required"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type UpdateCollectionResp struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}
