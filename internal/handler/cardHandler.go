package handler

import (
	"back/internal/exceptions"
	"back/internal/schemas"
	"back/internal/util"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *Handler) createCard(w http.ResponseWriter, r *http.Request) {
	var cardSchemaReq schemas.CreateCardReq

	if err := util.DecodeJSONRequest(r, &cardSchemaReq); err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidJSONFormat)
		return
	}

	if err := h.validator.ValidateWithDetailedErrors(&cardSchemaReq); err != nil {
		util.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	collectionIDStr := chi.URLParam(r, "collectionID")
	collectionID, err := strconv.Atoi(collectionIDStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidCollectionID)
		return
	}

	createdCard, err := h.services.CreateCard(&cardSchemaReq, collectionID)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusCreated, createdCard)

}

func (h *Handler) editCard(w http.ResponseWriter, r *http.Request) {
	var updatedCardSchemaReq schemas.UpdateCardReq

	if err := util.DecodeJSONRequest(r, &updatedCardSchemaReq); err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidJSONFormat)
		return
	}

	cardIDStr := chi.URLParam(r, "cardID")
	cardID, err := strconv.Atoi(cardIDStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidCardID)
		return
	}

	updatedCardSchemaReq.ID = cardID

	if err := h.validator.ValidateWithDetailedErrors(&updatedCardSchemaReq); err != nil {
		util.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	updatedCard, err := h.services.UpdateCard(&updatedCardSchemaReq)
	if err != nil {
		writeAppError(w, err)
		return
	}
	util.WriteJSON(w, http.StatusOK, updatedCard)
}

func (h *Handler) removeCard(w http.ResponseWriter, r *http.Request) {
	cardIDStr := chi.URLParam(r, "cardID")
	cardID, err := strconv.Atoi(cardIDStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidCardID)
		return
	}

	removeCardReq := schemas.RemoveCardReq{ID: cardID}

	if err := h.validator.ValidateWithDetailedErrors(&removeCardReq); err != nil {
		util.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = h.services.RemoveCard(&removeCardReq)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) getCardsByCollectionID(w http.ResponseWriter, r *http.Request) {
	collectionIDStr := chi.URLParam(r, "collectionID")
	collectionID, err := strconv.Atoi(collectionIDStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidCollectionID)
		return
	}

	cards, err := h.services.GetCardsByCollectionID(collectionID)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, cards)
}
