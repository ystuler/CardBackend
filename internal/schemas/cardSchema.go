package schemas

type CreateCardReq struct {
	Question string `json:"question" validate:"required"`
	Answer   string `json:"answer,omitempty"`
}

type CreateCardResp struct {
	ID           int    `json:"id"`
	Question     string `json:"question"`
	Answer       string `json:"answer"`
	CollectionID int    `json:"collection_id"`
}
