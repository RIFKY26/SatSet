package service

import (
	"satset2/rating/domain"
)

type RatingService struct {
	Repo domain.RatingRepository
}

func NewRatingService(r domain.RatingRepository) *RatingService {
	return &RatingService{Repo: r}
}

// FUNGSI INI SUDAH DIIMPLEMENTASIKAN AGAR UNIT TEST PASS
func (s *RatingService) SubmitRating(orderID, driverID string, score int) error {
	// Memanggil repository (di Unit Test akan memanggil Mock, di Functional memanggil DB asli)
	return s.Repo.SaveRating(orderID, driverID, score)
}
