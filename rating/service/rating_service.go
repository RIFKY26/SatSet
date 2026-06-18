package service

import (
	"errors"
	"satset2/rating/domain"
)

type RatingService struct {
	repo domain.RatingRepository
}

func NewRatingService(repo domain.RatingRepository) *RatingService {
	return &RatingService{repo: repo}
}

func (s *RatingService) SubmitRating(rating *domain.Rating) error {
	if rating.OrderID == "" || rating.DriverID == "" {
		return errors.New("order_id dan driver_id tidak boleh kosong")
	}
	if rating.Score < 1 || rating.Score > 5 {
		return errors.New("score harus antara 1 sampai 5")
	}

	return s.repo.SaveRating(rating)
}
