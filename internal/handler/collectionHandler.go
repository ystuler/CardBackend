package handler

import (
	"back/internal/exceptions"
	"back/internal/middleware"
	"back/internal/schemas"
	"back/internal/util"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

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

func (h *Handler) getCollectionByID(w http.ResponseWriter, r *http.Request) {
	collectionID, err := strconv.Atoi(chi.URLParam(r, "collectionID"))
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidCollectionID)
		return
	}

	collection, err := h.services.GetCollectionByID(collectionID)
	if err != nil {
		writeAppError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, collection)
}

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

func (h *Handler) removeCollection(w http.ResponseWriter, r *http.Request) {
	var removedCollectionSchema schemas.RemoveCollectionReq

	collectionID, err := strconv.Atoi(chi.URLParam(r, "collectionID"))
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidCollectionID)
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

func (h *Handler) startPractise(w http.ResponseWriter, r *http.Request) {
	var practiseSchemaReq schemas.TrainSchemaReq

	collectionID, err := strconv.Atoi(chi.URLParam(r, "collectionID"))
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, exceptions.ErrInvalidCollectionID)
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
