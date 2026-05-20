package service

import (
	"fmt"
	"satset2/promo-service/domain"
	"time"
)

type promoService struct {
	repo domain.PromoRepository
}

func NewPromoService(r domain.PromoRepository) domain.PromoService {
	return &promoService{repo: r}
}

func (s *promoService) ApplyPromo(code, userID string, orderValue float64, serviceType string) (*domain.PromoOutput, error) {
	promo, err := s.repo.FindByCode(code)
	if err != nil {
		return &domain.PromoOutput{IsValid: false, ErrorMessage: "Promo code not found"}, nil
	}

	// 1. Validasi Masa Berlaku [cite: 114]
	if time.Now().After(promo.ExpiryDate) {
		return &domain.PromoOutput{IsValid: false, ErrorMessage: "Promo has expired"}, nil
	}

	// 2. Validasi Kuota Umum [cite: 118]
	if promo.QuotaRemaining <= 0 {
		return &domain.PromoOutput{IsValid: false, ErrorMessage: "Promo quota is full"}, nil
	}

	// 3. Validasi Tipe Layanan
	if promo.ServiceType != serviceType {
		return &domain.PromoOutput{IsValid: false, ErrorMessage: fmt.Sprintf("Promo only valid for %s", promo.ServiceType)}, nil
	}

	// 4. Validasi Minimal Belanja
	if orderValue < promo.MinOrderValue {
		return &domain.PromoOutput{IsValid: false, ErrorMessage: "Minimum order value not met"}, nil
	}

	// 5. Validasi Penggunaan Per User
	usageCount, _ := s.repo.GetUserUsageCount(promo.PromoID, userID)
	if usageCount >= 1 {
		return &domain.PromoOutput{IsValid: false, ErrorMessage: "You have already used this promo"}, nil
	}

	// Perhitungan Diskon
	discount := orderValue * (promo.DiscountPercent / 100)
	if discount > promo.MaxDiscount {
		discount = promo.MaxDiscount
	}

	return &domain.PromoOutput{
		IsValid:        true,
		DiscountAmount: discount,
		FinalPrice:     orderValue - discount,
	}, nil
}
