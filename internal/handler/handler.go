package handler

import (
	"back/internal/service"
	"back/internal/util"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Handler struct {
	services  *service.Service
	validator *util.Validator
}

func NewHandler(services *service.Service, validator *util.Validator) *Handler {
	return &Handler{services: services, validator: validator}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("todo")) })

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", h.signUp)
		r.Post("/login", h.signIn)
	})

	return r
}
