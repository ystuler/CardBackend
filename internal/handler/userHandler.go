package handler

import (
	"back/internal/exceptions"
	"back/internal/middleware"
	"back/internal/schemas"
	"back/internal/util"
	"encoding/json"
	"net/http"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var userSchemaReq schemas.CreateUserReq

	//todo
	if err := json.NewDecoder(r.Body).Decode(&userSchemaReq); err != nil {
		http.Error(w, exceptions.ErrInvalidJSONFormat, http.StatusBadRequest)
		return
	}

	if err := h.validator.ValidateWithDetailedErrors(&userSchemaReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := h.services.SignUp(&userSchemaReq)
	if err != nil {
		http.Error(w, exceptions.ErrUserAlreadyExists, http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createdUser)
	if err != nil {
		http.Error(w, exceptions.ErrInternalServer, http.StatusInternalServerError)
		return
	}
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var userSchemaReq schemas.SignInReq

	if err := util.DecodeJSON(w, r, &userSchemaReq); err != nil {
		http.Error(w, exceptions.ErrInvalidJSONFormat, http.StatusBadRequest)
		return
	}

	if err := h.validator.ValidateWithDetailedErrors(&userSchemaReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.services.SignIn(&userSchemaReq)
	if err != nil {
		http.Error(w, exceptions.ErrInvalidCredentials, http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, exceptions.ErrInternalServer, http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		http.Error(w, exceptions.ErrInvalidToken, http.StatusUnauthorized)
		return
	}

	resp, err := h.services.GetProfile(userID)
	if err != nil {
		http.Error(w, exceptions.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, exceptions.ErrInternalServer, http.StatusInternalServerError)
		return
	}
}

func (h *Handler) updateUsername(w http.ResponseWriter, r *http.Request) {
	var updateUsernameReq schemas.UpdateUsernameReq

	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		http.Error(w, exceptions.ErrInvalidToken, http.StatusUnauthorized)
		return
	}

	updateUsernameReq.ID = userID

	if err := util.DecodeJSON(w, r, &updateUsernameReq); err != nil {
		http.Error(w, exceptions.ErrInvalidJSONFormat, http.StatusBadRequest)
		return
	}

	if err := h.validator.ValidateWithDetailedErrors(&updateUsernameReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.services.UpdateUsername(&updateUsernameReq)
	if err != nil {
		http.Error(w, exceptions.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, exceptions.ErrInternalServer, http.StatusInternalServerError)
	}
}

func (h *Handler) updatePassword(w http.ResponseWriter, r *http.Request) {
	var updatePasswordReq schemas.UpdatePasswordReq
	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		http.Error(w, exceptions.ErrInvalidToken, http.StatusUnauthorized)
		return
	}

	updatePasswordReq.ID = userID

	if err := util.DecodeJSON(w, r, &updatePasswordReq); err != nil {
		http.Error(w, exceptions.ErrInvalidJSONFormat, http.StatusBadRequest)
		return
	}
	if err := h.validator.ValidateWithDetailedErrors(&updatePasswordReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.services.UpdatePassword(&updatePasswordReq); err != nil {
		http.Error(w, exceptions.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
