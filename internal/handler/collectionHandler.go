package handler

import (
	"back/internal/middleware"
	"back/internal/schemas"
	"back/internal/util"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *Handler) createCollection(w http.ResponseWriter, r *http.Request) {
	var collectionSchemaReq schemas.CreateCollectionReq

	if err := util.DecodeJSON(w, r, &collectionSchemaReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
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

func (h *Handler) editCollection(w http.ResponseWriter, r *http.Request) {
	var updatedCollectionSchema schemas.UpdateCollectionReq

	if err := util.DecodeJSON(w, r, &updatedCollectionSchema); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collectionID, err := strconv.Atoi(chi.URLParam(r, "collectionID"))
	if err != nil {
		http.Error(w, "id must be an integer", http.StatusBadRequest)
	}
	updatedCollectionSchema.ID = collectionID

	if err := h.validator.Validate(&updatedCollectionSchema); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	updatedCollection, err := h.services.UpdateCollection(&updatedCollectionSchema, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(updatedCollection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) removeCollection(w http.ResponseWriter, r *http.Request) {
	var removedCollectionSchema schemas.RemoveCollectionReq

	if err := util.DecodeJSON(w, r, &removedCollectionSchema); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collectionID, err := strconv.Atoi(chi.URLParam(r, "collectionID"))
	if err != nil {
		http.Error(w, "id must be an integer", http.StatusBadRequest)
	}
	removedCollectionSchema.ID = collectionID

	if err := h.validator.Validate(&removedCollectionSchema); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.services.RemoveCollection(&removedCollectionSchema)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
