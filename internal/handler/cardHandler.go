package handler

import (
	"back/internal/schemas"
	"back/internal/util"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *Handler) createCard(w http.ResponseWriter, r *http.Request) {
	var cardSchemaReq schemas.CreateCardReq

	if err := util.DecodeJSON(w, r, &cardSchemaReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Validate(&cardSchemaReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collectionIDStr := chi.URLParam(r, "collectionID")
	collectionID, err := strconv.Atoi(collectionIDStr)
	if err != nil {
		http.Error(w, "Invalid collection ID", http.StatusBadRequest)
		return
	}

	createdCard, err := h.services.CreateCard(&cardSchemaReq, collectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdCard); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
