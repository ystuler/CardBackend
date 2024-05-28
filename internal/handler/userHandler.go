package handler

import (
	"back/internal/schemas"
	"back/internal/util"
	"encoding/json"
	"net/http"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	var userSchemaReq schemas.CreateUserReq

	if err := json.NewDecoder(r.Body).Decode(&userSchemaReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Validate(&userSchemaReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := h.services.CreateUser(&userSchemaReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createdUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	var userSchemaReq schemas.SignInReq

	if err := util.DecodeJSON(w, r, &userSchemaReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Validate(&userSchemaReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.services.SignIn(&userSchemaReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
