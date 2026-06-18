package service

import (
	"satset2/promo_service/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 1. Buat Mock Repository Baru
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

func (m *MockPromoRepo) UpdateQuota(promo *domain.Promo) error {
	args := m.Called(promo)
	return args.Error(0)
}

// 2. Test Skenario Sukses
func TestApplyPromoSuccess(t *testing.T) {
	mockRepo := new(MockPromoRepo)

	// Data promo bohongan untuk di-test
	dummyPromo := &domain.Promo{
		PromoID:       "P-001",
		PromoCode:     "SATSET50",
		MinOrderValue: 20000,
		MaxDiscount:   15000,
		DiscountPct:   50,
		Quota:         10,
	}

	// Atur agar mock mengembalikan dummyPromo
	mockRepo.On("FindByCode", "SATSET50").Return(dummyPromo, nil)
	mockRepo.On("UpdateQuota", dummyPromo).Return(nil)

	// Jalankan service-nya
	svc := NewPromoService(mockRepo)

	// Kita tes beli seharga Rp 40.000
	discount, err := svc.ApplyPromo("SATSET50", 40000)

	// Validasi hasil
	assert.NoError(t, err)
	// Diskon 50% dari 40.000 adalah 20.000, tapi karena MaxDiscount 15.000, maka hasil harus 15.000
	assert.Equal(t, float64(15000), discount)
}
