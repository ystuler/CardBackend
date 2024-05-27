package handler

import (
	"back/internal/schemas"
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

	userSchemaResp := schemas.CreateUserResp{
		ID:       createdUser.ID,
		Username: createdUser.Username,
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(userSchemaResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
