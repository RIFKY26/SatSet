package handler

import (
	"encoding/json"
	"net/http"

	"satset2/order-service/domain"
	"satset2/order-service/repository"
	"satset2/order-service/service"
)

// CreateOrderHandler menangani permintaan HTTP
func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Gunakan OrderRequest dari package domain
	var req domain.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Panggil logika bisnis dari package service, dan gunakan repository yang baru
	repo := repository.DefaultOrderRepository{}
	order, err := service.CreateOrder(req, repo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}
