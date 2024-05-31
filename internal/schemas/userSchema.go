package schemas

type CreateUserReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CreateUserResp struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

type SignInReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignInResp struct {
	Token    string `json:"token"`
	ID       int    `json:"id"`
	Username string `json:"username"`
}
