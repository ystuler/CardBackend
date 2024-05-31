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

func (h *Handler) editCard(w http.ResponseWriter, r *http.Request) {
	var updatedCardSchemaReq schemas.UpdateCardReq

	if err := util.DecodeJSON(w, r, &updatedCardSchemaReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cardIDStr := chi.URLParam(r, "cardID")
	cardID, err := strconv.Atoi(cardIDStr)
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	updatedCardSchemaReq.ID = cardID

	if err := h.validator.Validate(&updatedCardSchemaReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedCard, err := h.services.UpdateCard(&updatedCardSchemaReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(updatedCard); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) removeCard(w http.ResponseWriter, r *http.Request) {
	cardIDStr := chi.URLParam(r, "cardID")
	cardID, err := strconv.Atoi(cardIDStr)
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	removeCardReq := schemas.RemoveCardReq{ID: cardID}

	if err := h.validator.Validate(&removeCardReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.services.RemoveCard(&removeCardReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Card successfully removed"))
}

func (h *Handler) getCardsByCollectionID(w http.ResponseWriter, r *http.Request) {

	collectionIDStr := chi.URLParam(r, "collectionID")
	collectionID, err := strconv.Atoi(collectionIDStr)
	if err != nil {
		http.Error(w, "Invalid collection ID", http.StatusBadRequest)
		return
	}

	cards, err := h.services.GetCardsByCollectionID(collectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(cards); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
