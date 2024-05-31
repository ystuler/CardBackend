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

type Profile struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type GetProfileReq struct {
	ID int `validate:"required,gt=0"`
}

type GetProfileResp struct {
	Profile Profile `json:"profile"`
}

type UpdateUsernameReq struct {
	ID       int    `validate:"required,gt=0"`
	Username string `json:"username" validate:"required"`
}

type UpdateUsernameResp struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type UpdatePasswordReq struct {
	ID          int    `validate:"required,gt=0"`
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}
