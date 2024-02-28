package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Hymiside/hezzl-api/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type service interface {
	Goods(ctx context.Context, limit, offset int) (models.GoodsResponse, error)
	CreateGood(ctx context.Context, projectID int, name string) (models.Good, error)
	UpdateGood(ctx context.Context, good models.Good, goodID, projectID int) (models.Good, error)
}

type Handler struct {
	service  service
	validate *validator.Validate
}

func NewHandler(service service) *Handler {
	return &Handler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *Handler) NewRoutes() *chi.Mux {
	mux := chi.NewRouter()
	mux.Route("/api/", func(r chi.Router) {
		r.Get("/goods", h.Goods)
		r.Post("/goods", h.CreateGoods)
		r.Patch("/goods", h.UpdateGoods)
		r.Delete("/goods", h.DeleteGoods)
	})
	return mux
}

func (h *Handler) Goods(w http.ResponseWriter, r *http.Request) {
	limit, offset := r.URL.Query().Get("limit"), r.URL.Query().Get("offset")

	var (
		limitInt, offsetInt = 10, 0
		err                 error
	)

	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			http.Error(w, "limit must be an integer", http.StatusBadRequest)
			return
		}
	}
	if offset != "" {
		offsetInt, err = strconv.Atoi(offset)
		if err != nil {
			http.Error(w, "offset must be an integer", http.StatusBadRequest)
			return
		}
	}

	goodsResponse, err := h.service.Goods(r.Context(), limitInt, offsetInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(goodsResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreateGoods(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("projectId")
	if projectID == "" {
		http.Error(w, "projectId is required", http.StatusBadRequest)
		return
	}

	projectIDInt, err := strconv.Atoi(projectID)
	if err != nil {
		http.Error(w, "projectId must be an integer", http.StatusBadRequest)
		return
	}

	var data models.CreateGoodRequest
	if err = json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.validate.Struct(data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	good, err := h.service.CreateGood(r.Context(), projectIDInt, data.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(good); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateGoods(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("projectId")
	if projectID == "" {
		http.Error(w, "projectId is required", http.StatusBadRequest)
		return
	}

	goodID := r.URL.Query().Get("id")
	if goodID == "" {
		http.Error(w, "goodId is required", http.StatusBadRequest)
		return
	}

	projectIDInt, err := strconv.Atoi(projectID)
	if err != nil {
		http.Error(w, "projectId must be an integer", http.StatusBadRequest)
		return
	}

	goodIDInt, err := strconv.Atoi(goodID)
	if err != nil {
		http.Error(w, "goodId must be an integer", http.StatusBadRequest)
		return
	}

	var data models.Good
	if err = json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	goodUpdated, err := h.service.UpdateGood(r.Context(), data, goodIDInt, projectIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(goodUpdated); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteGoods(w http.ResponseWriter, r *http.Request) {}
