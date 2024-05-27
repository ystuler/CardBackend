package schemas

type CreateUserReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CreateUserResp struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}
