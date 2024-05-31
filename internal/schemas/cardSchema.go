package schemas

type CreateCardReq struct {
	Question string `json:"question" validate:"required"`
	Answer   string `json:"answer,omitempty"`
}

type CreateCardResp struct {
	ID           int    `json:"id"`
	Question     string `json:"question"`
	Answer       string `json:"answer"`
	CollectionID int    `json:"collectionID"`
}

type UpdateCardReq struct {
	ID       int     `validate:"required,gt=0"`
	Question *string `json:"question,omitempty"`
	Answer   *string `json:"answer,omitempty"`
}

type UpdateCardResp struct {
	ID           int    `json:"id"`
	Question     string `json:"question"`
	Answer       string `json:"answer"`
	CollectionID int    `json:"collection_id"`
}

type RemoveCardReq struct {
	ID int `validate:"required,gt=0"`
}

type GetCardsByCollectionIDReq struct {
	CollectionID int `validate:"required,gt=0"`
}

type CardsByCollectionID struct {
	ID       int    `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type GetCardByCollectionIDResp struct {
	Cards []CardsByCollectionID `json:"cards"`
}
