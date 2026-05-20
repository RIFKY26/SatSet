package repository

import (
	"errors"
	"satset2/promo-service/domain"
)

type promoRepository struct {
	// Di sini biasanya ada *sql.DB atau library database lainnya
}

func NewPromoRepository() domain.PromoRepository {
	return &promoRepository{}
}

func (r *promoRepository) FindByCode(code string) (*domain.Promo, error) {
	// Logika pencarian data ke database
	return nil, errors.New("not implemented")
}

func (r *promoRepository) GetUserUsageCount(promoID, userID string) (int, error) {
	// Menghitung penggunaan user
	return 0, nil
}

func (r *promoRepository) UpdateQuota(promoID string, amount int) error {
	// Update kuota di DB
	return nil
}

func (r *promoRepository) RecordUsage(history *domain.UsageHistory) error {
	// Simpan history penggunaan
	return nil
}
