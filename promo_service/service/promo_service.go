package service

import (
	"errors"
	"satset2/promo_service/domain"
)

type PromoService struct {
	repo domain.PromoRepository
}

func NewPromoService(repo domain.PromoRepository) *PromoService {
	return &PromoService{repo: repo}
}

func (s *PromoService) ApplyPromo(code string, orderValue float64) (float64, error) {
	// 1. Cari promo di database
	promo, err := s.repo.FindByCode(code)
	if err != nil {
		return 0, errors.New("kode promo tidak ditemukan atau tidak valid")
	}

	// 2. Cek validasi kuota dan syarat minimum
	if promo.Quota <= 0 {
		return 0, errors.New("maaf, kuota promo ini sudah habis")
	}
	if orderValue < promo.MinOrderValue {
		return 0, errors.New("minimum transaksi tidak terpenuhi untuk promo ini")
	}

	// 3. Hitung diskon
	discount := orderValue * (promo.DiscountPct / 100)
	if discount > promo.MaxDiscount {
		discount = promo.MaxDiscount
	}

	// 4. Kurangi kuota promo di database agar tidak dipakai terus-terusan
	promo.Quota -= 1
	s.repo.UpdateQuota(promo)

	return discount, nil
}
