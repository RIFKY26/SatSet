package handler

import (
	"encoding/json"
	"net/http"
	"satset2/promo-service/domain"
)

type PromoHandler struct {
	service domain.PromoService
}

func NewPromoHandler(s domain.PromoService) *PromoHandler {
	return &PromoHandler{service: s}
}

func (h *PromoHandler) HandleApplyPromo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req domain.PromoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.service.ApplyPromo(req.PromoCode, req.UserID, req.OrderValue, req.ServiceType)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if !result.IsValid {
		w.WriteHeader(http.StatusUnprocessableEntity)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(result)
}
