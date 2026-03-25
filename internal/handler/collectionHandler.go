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

// createCollection creates a new collection for the current user.
// @Summary Create collection
// @Tags collections
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body schemas.CreateCollectionReq true "collection payload"
// @Success 201 {object} schemas.CreateCollectionResp
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 422 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /collections/ [post]
func (h *Handler) createCollection(w http.ResponseWriter, r *http.Request) {
	var collectionSchemaReq schemas.CreateCollectionReq

	if err := util.DecodeJSONRequest(r, &collectionSchemaReq); err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidJSONFormat)
		return
	}

	if err := h.validator.ValidateWithDetailedErrors(&collectionSchemaReq); err != nil {
		util.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
		return
	}

	createdCollection, err := h.services.CreateCollection(&collectionSchemaReq, userID)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusCreated, createdCollection)
}

// getCollectionByID returns collection by ID.
// @Summary Get collection by ID
// @Tags collections
// @Security BearerAuth
// @Produce json
// @Param collectionID path int true "collection id"
// @Success 200 {object} schemas.GetCollectionByIDResp
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Router /collections/{collectionID}/ [get]
func (h *Handler) getCollectionByID(w http.ResponseWriter, r *http.Request) {
	collectionID, err := strconv.Atoi(chi.URLParam(r, "collectionID"))
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

	collection, err := h.services.GetCollectionByID(collectionID)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, collection)
}

// editCollection fully updates collection.
// @Summary Full update collection
// @Tags collections
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param collectionID path int true "collection id"
// @Param request body schemas.UpdateCollectionReq true "collection payload"
// @Success 200 {object} schemas.UpdateCollectionResp
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 422 {object} util.ErrorResponse
// @Router /collections/{collectionID}/ [put]
func (h *Handler) editCollection(w http.ResponseWriter, r *http.Request) {
	var updatedCollectionSchema schemas.UpdateCollectionReq

	if err := util.DecodeJSONRequest(r, &updatedCollectionSchema); err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidJSONFormat)
		return
	}

	collectionID, err := strconv.Atoi(chi.URLParam(r, "collectionID"))
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

	updatedCollectionSchema.ID = collectionID

	if err := h.validator.ValidateWithDetailedErrors(&updatedCollectionSchema); err != nil {
		util.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	updatedCollection, err := h.services.UpdateCollection(&updatedCollectionSchema)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, updatedCollection)
}

// patchCollection partially updates collection.
// @Summary Partial update collection
// @Tags collections
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param collectionID path int true "collection id"
// @Param request body schemas.PatchCollectionReq true "collection payload"
// @Success 200 {object} schemas.UpdateCollectionResp
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 422 {object} util.ErrorResponse
// @Router /collections/{collectionID}/ [patch]
func (h *Handler) patchCollection(w http.ResponseWriter, r *http.Request) {
	var patchCollectionSchema schemas.PatchCollectionReq

	if err := util.DecodeJSONRequest(r, &patchCollectionSchema); err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidJSONFormat)
		return
	}

	collectionID, err := strconv.Atoi(chi.URLParam(r, "collectionID"))
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

	patchCollectionSchema.ID = collectionID

	if patchCollectionSchema.Name == nil && patchCollectionSchema.Description == nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidRequestBody)
		return
	}

	if err := h.validator.ValidateWithDetailedErrors(&patchCollectionSchema); err != nil {
		util.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	updatedCollection, err := h.services.PatchCollection(&patchCollectionSchema)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, updatedCollection)
}

// removeCollection deletes collection.
// @Summary Delete collection
// @Tags collections
// @Security BearerAuth
// @Produce json
// @Param collectionID path int true "collection id"
// @Success 204
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Router /collections/{collectionID}/ [delete]
func (h *Handler) removeCollection(w http.ResponseWriter, r *http.Request) {
	var removedCollectionSchema schemas.RemoveCollectionReq

	collectionID, err := strconv.Atoi(chi.URLParam(r, "collectionID"))
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

	removedCollectionSchema.ID = collectionID

	if err := h.validator.ValidateWithDetailedErrors(&removedCollectionSchema); err != nil {
		util.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = h.services.RemoveCollection(&removedCollectionSchema)
	if err != nil {
		writeAppError(w, err)
		return
	}
	util.WriteJSON(w, http.StatusNoContent, nil)
}

// getAllCollections returns all current user collections.
// @Summary Get all collections
// @Tags collections
// @Security BearerAuth
// @Produce json
// @Success 200 {object} schemas.AllCollectionsResp
// @Failure 401 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /collections/ [get]
func (h *Handler) getAllCollections(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusUnauthorized, exceptions.ErrInvalidToken)
		return
	}

	allCollections, err := h.services.GetAllCollections(userID)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, allCollections)

}

// startPractise returns shuffled cards for training.
// @Summary Train cards
// @Tags collections
// @Security BearerAuth
// @Produce json
// @Param collectionID path int true "collection id"
// @Success 200 {object} schemas.TrainSchemaResp
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 422 {object} util.ErrorResponse
// @Router /collections/{collectionID}/train [get]
func (h *Handler) startPractise(w http.ResponseWriter, r *http.Request) {
	var practiseSchemaReq schemas.TrainSchemaReq

	collectionID, err := strconv.Atoi(chi.URLParam(r, "collectionID"))
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

	practiseSchemaReq.ID = collectionID

	if err := h.validator.ValidateWithDetailedErrors(&practiseSchemaReq); err != nil {
		util.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	randomCards, err := h.services.TrainCards(&practiseSchemaReq)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, randomCards)
}
