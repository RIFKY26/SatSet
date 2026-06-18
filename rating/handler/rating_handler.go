package handler

import (
	"encoding/json"
	"net/http"

	"satset2/rating/domain"
	"satset2/rating/service"
)

type RatingHandler struct {
	ratingService *service.RatingService
}

func NewRatingHandler(s *service.RatingService) *RatingHandler {
	return &RatingHandler{ratingService: s}
}

func (h *RatingHandler) SubmitRatingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req domain.Rating
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := h.ratingService.SubmitRating(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Terima kasih atas penilaian Anda!",
	})
}
