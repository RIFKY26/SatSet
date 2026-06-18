package handler

import (
	"encoding/json"
	"net/http"

	"satset2/matching-service/domain"
	"satset2/matching-service/service"
)

type MatchHandler struct {
	matchService *service.MatchService
}

func NewMatchHandler(s *service.MatchService) *MatchHandler {
	return &MatchHandler{matchService: s}
}

func (h *MatchHandler) MatchDriverHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req domain.MatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	result, err := h.matchService.MatchDriver(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}