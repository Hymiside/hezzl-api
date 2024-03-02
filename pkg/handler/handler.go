package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Hymiside/hezzl-api/pkg/custerrors"
	"github.com/Hymiside/hezzl-api/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type service interface {
	Goods(ctx context.Context, limit, offset int) (models.GoodsResponse, error)
	CreateGood(ctx context.Context, projectID int, name string) (models.Good, error)
	UpdateGood(ctx context.Context, good models.Good, goodID, projectID int) (models.Good, error)
	DeleteGood(ctx context.Context, goodID, projectID int) (models.Good, error)
	ReprioritizeGood(ctx context.Context, goodID, projectID, priority int) ([]models.ReprioritizeGoodResponse, error)
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
	mux.Route("/goods/", func(r chi.Router) {
		r.Get("/list", h.Goods)
		r.Post("/create", h.CreateGood)
		r.Patch("/update", h.UpdateGood)
		r.Delete("/delete", h.DeleteGood)
		r.Patch("/reprioritize", h.ReprioritizeGood)
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

func (h *Handler) CreateGood(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) UpdateGood(w http.ResponseWriter, r *http.Request) {
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
		if errors.Is(err, custerrors.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
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

func (h *Handler) DeleteGood(w http.ResponseWriter, r *http.Request) {
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

	good, err := h.service.DeleteGood(r.Context(), goodIDInt, projectIDInt)
	if err != nil {
		if errors.Is(err, custerrors.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(
		map[string]interface{}{
			"id": good.ID,
			"projectId": good.ProjectID,
			"removed": good.Removed,
		},
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ReprioritizeGood(w http.ResponseWriter, r *http.Request) {
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

	var data models.ReprioritizeGoodRequest
	if err = json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.validate.Struct(data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reprioritizedGood, err := h.service.ReprioritizeGood(r.Context(), goodIDInt, projectIDInt, data.Priority)
	if err != nil {
		if errors.Is(err, custerrors.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(
		map[string]interface{}{"priorities": reprioritizedGood},
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
