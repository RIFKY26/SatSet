package repository

import (
	"satset2/rating/domain"

	"gorm.io/gorm"
)

type SqlRatingRepository struct {
	DB *gorm.DB
}

func NewSqlRatingRepository(db *gorm.DB) domain.RatingRepository {
	return &SqlRatingRepository{DB: db}
}

func (r *SqlRatingRepository) SaveRating(rating *domain.Rating) error {
	return r.DB.Create(rating).Error
}
