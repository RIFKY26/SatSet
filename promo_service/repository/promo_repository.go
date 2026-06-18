package repository

import (
	"satset2/promo_service/domain"

	"gorm.io/gorm"
)

type SqlPromoRepository struct {
	DB *gorm.DB
}

func NewSqlPromoRepository(db *gorm.DB) domain.PromoRepository {
	return &SqlPromoRepository{DB: db}
}

func (r *SqlPromoRepository) FindByCode(code string) (*domain.Promo, error) {
	var promo domain.Promo
	err := r.DB.Where("promo_code = ?", code).First(&promo).Error
	return &promo, err
}

func (r *SqlPromoRepository) UpdateQuota(promo *domain.Promo) error {
	return r.DB.Save(promo).Error
}
