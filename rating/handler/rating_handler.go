package handler

import (
	"net/http"
	"satset2/rating/service"
)

type RatingHandler struct {
	RatingService *service.RatingService
}

func NewRatingHandler(svc *service.RatingService) *RatingHandler {
	return &RatingHandler{RatingService: svc}
}

func (h *RatingHandler) HandleSubmitRating(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"message": "TDD: handler HTTP rating belum diimplementasikan"}`))
}
