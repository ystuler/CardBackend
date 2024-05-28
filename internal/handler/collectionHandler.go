package handler

import (
	"back/internal/middleware"
	"back/internal/schemas"
	"encoding/json"
	"net/http"
)

func (h *Handler) createCollection(w http.ResponseWriter, r *http.Request) {
	var collectionSchemaReq schemas.CreateCollectionReq

	if err := json.NewDecoder(r.Body).Decode(&collectionSchemaReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := h.validator.Validate(&collectionSchemaReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	createdCollection, err := h.services.CreateCollection(&collectionSchemaReq, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createdCollection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
