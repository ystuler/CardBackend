package handler

import (
	"back/internal/exceptions"
	"back/internal/middleware"
	"back/internal/schemas"
	"back/internal/util"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// createCard creates a card in collection.
// @Summary Create card
// @Tags cards
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param collectionID path int true "collection id"
// @Param request body schemas.CreateCardReq true "card payload"
// @Success 201 {object} schemas.CreateCardResp
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 422 {object} util.ErrorResponse
// @Router /collections/{collectionID}/cards/ [post]
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

	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
		return
	}

	if err := h.services.EnsureCollectionAccess(collectionID, userID); err != nil {
		writeAppError(w, err)
		return
	}

	createdCard, err := h.services.CreateCard(&cardSchemaReq, collectionID)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusCreated, createdCard)

}

// editCard updates card (used for PUT and PATCH routes).
// @Summary Update card
// @Tags cards
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param cardID path int true "card id"
// @Param request body schemas.UpdateCardReq true "card payload"
// @Success 200 {object} schemas.UpdateCardResp
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 422 {object} util.ErrorResponse
// @Router /cards/{cardID}/ [put]
// @Router /cards/{cardID}/ [patch]
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

	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
		return
	}

	if err := h.services.EnsureCardAccess(cardID, userID); err != nil {
		writeAppError(w, err)
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

// removeCard deletes card.
// @Summary Delete card
// @Tags cards
// @Security BearerAuth
// @Produce json
// @Param cardID path int true "card id"
// @Success 204
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Router /cards/{cardID}/ [delete]
func (h *Handler) removeCard(w http.ResponseWriter, r *http.Request) {
	cardIDStr := chi.URLParam(r, "cardID")
	cardID, err := strconv.Atoi(cardIDStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidCardID)
		return
	}

	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
		return
	}

	if err := h.services.EnsureCardAccess(cardID, userID); err != nil {
		writeAppError(w, err)
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

// getCardsByCollectionID returns all cards from collection.
// @Summary Get cards by collection
// @Tags cards
// @Security BearerAuth
// @Produce json
// @Param collectionID path int true "collection id"
// @Success 200 {object} schemas.GetCardByCollectionIDResp
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Router /collections/{collectionID}/cards/ [get]
func (h *Handler) getCardsByCollectionID(w http.ResponseWriter, r *http.Request) {
	collectionIDStr := chi.URLParam(r, "collectionID")
	collectionID, err := strconv.Atoi(collectionIDStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidCollectionID)
		return
	}

	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
		return
	}

	if err := h.services.EnsureCollectionAccess(collectionID, userID); err != nil {
		writeAppError(w, err)
		return
	}

	cards, err := h.services.GetCardsByCollectionID(collectionID)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, cards)
}
