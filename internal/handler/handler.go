package handler

import (
	"back/internal/service"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("todo")) })
	r.Get("/user", func(w http.ResponseWriter, r *http.Request) {

	})
	return r
}
