package handlers

import (
	"net/http"
	"wb_order_service/internal/cache"
	"wb_order_service/internal/database"
	"wb_order_service/pkg/response"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	cache    *cache.Cache
	database database.OrderStore
}

func NewHandler(cache *cache.Cache, database database.OrderStore) *Handler {
	return &Handler{
		cache:    cache,
		database: database,
	}
}

func (h *Handler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	orderUID := chi.URLParam(r, "order_uid")
	if orderUID == "" {
		response.Error(w, http.StatusBadRequest, "order_uid is required")
		return
	}

	if order, exists := h.cache.Get(orderUID); exists {
		response.Json(w, http.StatusOK, order)
		return
	}

	order, err := h.database.GetOrderByID(r.Context(), orderUID)
	if err != nil {
		if err.Error() == "order not found" {
			response.Error(w, http.StatusNotFound, "order not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.cache.Set(order)
	response.Json(w, http.StatusOK, order)
}
