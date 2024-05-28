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
