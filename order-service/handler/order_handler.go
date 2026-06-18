package handler

import (
	"encoding/json"
	"net/http"

	"satset2/order-service/domain"
	"satset2/order-service/service"
)

type OrderHandler struct {
	Repo domain.OrderRepository
}

func NewOrderHandler(repo domain.OrderRepository) *OrderHandler {
	return &OrderHandler{Repo: repo}
}

func (h *OrderHandler) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req domain.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Panggil logika bisnis dengan mengirimkan repo yang asli
	order, err := service.CreateOrder(req, h.Repo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}