package repository

import (
	"errors"
	"satset2/rating/domain"
)

type SqlRatingRepository struct {
	// Nanti diisi koneksi database DB
}

func NewRatingRepository() domain.RatingRepository {
	return &SqlRatingRepository{}
}

func (r *SqlRatingRepository) SaveRating(orderID string, driverID string, score int) error {
	return errors.New("TDD: repository SaveRating belum diimplementasikan")
}
