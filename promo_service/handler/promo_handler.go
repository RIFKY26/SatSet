package handler

import (
	"encoding/json"
	"net/http"

	"satset2/promo_service/service"
)

type PromoHandler struct {
	promoService *service.PromoService
}

func NewPromoHandler(s *service.PromoService) *PromoHandler {
	return &PromoHandler{promoService: s}
}

func (h *PromoHandler) ApplyPromoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Code       string  `json:"promo_code"`
		OrderValue float64 `json:"order_value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	discount, err := h.promoService.ApplyPromo(req.Code, req.OrderValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"promo_code": req.Code,
		"discount":   discount,
		"status":     "applied",
	})
}
