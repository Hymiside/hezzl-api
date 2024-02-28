package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type service interface {}

type Handler struct {
	service service
	validate *validator.Validate
}

func NewHandler(service service) *Handler {
	return &Handler{
		service: service,
		validate: validator.New(),
	}
}

func (h *Handler) NewRoutes() *chi.Mux {
	return chi.NewRouter()
}