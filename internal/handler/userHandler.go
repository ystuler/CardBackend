package handler

import (
	"back/internal/exceptions"
	"back/internal/middleware"
	"back/internal/schemas"
	"back/internal/util"
	"net/http"
)

// SignUp registers a new user.
// @Summary Register user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body schemas.CreateUserReq true "signup payload"
// @Success 201 {object} schemas.CreateUserResp
// @Failure 400 {object} util.ErrorResponse
// @Failure 422 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /auth/signup [post]
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var userSchemaReq schemas.CreateUserReq

	if err := util.DecodeJSONRequest(r, &userSchemaReq); err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidJSONFormat)
		return
	}

	if err := h.validator.ValidateWithDetailedErrors(&userSchemaReq); err != nil {
		util.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	createdUser, err := h.services.SignUp(&userSchemaReq)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusCreated, createdUser)
}

// SignIn authenticates user and returns JWT.
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body schemas.SignInReq true "signin payload"
// @Success 200 {object} schemas.SignInResp
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 422 {object} util.ErrorResponse
// @Router /auth/login [post]
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var userSchemaReq schemas.SignInReq

	if err := util.DecodeJSONRequest(r, &userSchemaReq); err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidJSONFormat)
		return
	}

	if err := h.validator.ValidateWithDetailedErrors(&userSchemaReq); err != nil {
		util.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	resp, err := h.services.SignIn(&userSchemaReq)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, resp)
}

// getProfile returns current user profile.
// @Summary Get profile
// @Tags profile
// @Security BearerAuth
// @Produce json
// @Success 200 {object} schemas.GetProfileResp
// @Failure 401 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /profile/ [get]
func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
		return
	}

	resp, err := h.services.GetProfile(userID)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, resp)
}

// updateUsername updates current user username.
// @Summary Update username
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body schemas.UpdateUsernameBody true "username payload"
// @Success 200 {object} schemas.UpdateUsernameResp
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 422 {object} util.ErrorResponse
// @Router /profile/username [put]
func (h *Handler) updateUsername(w http.ResponseWriter, r *http.Request) {
	var updateUsernameBody schemas.UpdateUsernameBody

	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
		return
	}

	if err := util.DecodeJSONRequest(r, &updateUsernameBody); err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidJSONFormat)
		return
	}

	updateUsernameReq := schemas.UpdateUsernameReq{
		ID:       userID,
		Username: updateUsernameBody.Username,
	}

	if err := h.validator.ValidateWithDetailedErrors(&updateUsernameReq); err != nil {
		util.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	resp, err := h.services.UpdateUsername(&updateUsernameReq)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, resp)
}

// updatePassword updates current user password.
// @Summary Update password
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body schemas.UpdatePasswordReq true "password payload"
// @Success 204
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 422 {object} util.ErrorResponse
// @Router /profile/password [put]
func (h *Handler) updatePassword(w http.ResponseWriter, r *http.Request) {
	var updatePasswordReq schemas.UpdatePasswordReq
	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
		return
	}

	updatePasswordReq.ID = userID

	if err := util.DecodeJSONRequest(r, &updatePasswordReq); err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidJSONFormat)
		return
	}
	if err := h.validator.ValidateWithDetailedErrors(&updatePasswordReq); err != nil {
		util.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	if err := h.services.UpdatePassword(&updatePasswordReq); err != nil {
		writeAppError(w, err)
		return
	}
	util.WriteJSON(w, http.StatusNoContent, nil)
}

// LogOut is a formal logout endpoint for JWT architecture.
// @Summary Logout user
// @Tags auth
// @Produce json
// @Success 200 {object} schemas.LogOutResp
// @Router /auth/logout [post]
func (h *Handler) LogOut(w http.ResponseWriter, r *http.Request) {
	util.WriteJSON(w, http.StatusOK, schemas.LogOutResp{Message: "logged out"})
}
