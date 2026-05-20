package service

import (
	"satset2/promo-service/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPromoRepo struct {
	mock.Mock
}

func (m *MockPromoRepo) FindByCode(code string) (*domain.Promo, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Promo), args.Error(1)
}

func (m *MockPromoRepo) GetUserUsageCount(pID, uID string) (int, error) {
	args := m.Called(pID, uID)
	return args.Int(0), args.Error(1)
}
func (m *MockPromoRepo) UpdateQuota(id string, amt int) error     { return nil }
func (m *MockPromoRepo) RecordUsage(h *domain.UsageHistory) error { return nil }

func TestApplyPromoSuccess(t *testing.T) {
	mockRepo := new(MockPromoRepo)
	now := time.Now().Add(24 * time.Hour)

	dummyPromo := &domain.Promo{
		PromoID: "P1", PromoCode: "HEMAT", MinOrderValue: 10000,
		MaxDiscount: 5000, DiscountPercent: 50, QuotaRemaining: 10,
		ExpiryDate: now, ServiceType: "motor",
	}
	mockRepo.On("FindByCode", "HEMAT").Return(dummyPromo, nil)
	mockRepo.On("GetUserUsageCount", "P1", "USER1").Return(0, nil)

	svc := NewPromoService(mockRepo)
	res, _ := svc.ApplyPromo("HEMAT", "USER1", 20000, "motor")

	assert.True(t, res.IsValid)
	assert.Equal(t, float64(5000), res.DiscountAmount)
	assert.Equal(t, float64(15000), res.FinalPrice)
}
