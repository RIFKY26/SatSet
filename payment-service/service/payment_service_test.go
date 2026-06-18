package service

import (
	"satset2/payment-service/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPaymentRepo struct {
	mock.Mock
}

func (m *MockPaymentRepo) CreateTransaction(tx *domain.Transaction) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockPaymentRepo) FindByIdempotencyKey(key string) (*domain.Transaction, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

func TestProcessPaymentSuccess(t *testing.T) {
	mockRepo := new(MockPaymentRepo)
	// Pura-puranya transaksi ini belum pernah terjadi
	mockRepo.On("FindByIdempotencyKey", "KEY123").Return(nil, nil)
	mockRepo.On("CreateTransaction", mock.Anything).Return(nil)

	svc := NewPaymentService(mockRepo)

	// Format request yang baru
	req := PaymentRequest{
		OrderID:        "ORD-1",
		UserID:         1,
		Amount:         50000,
		PaymentMethod:  "gopay",
		IdempotencyKey: "KEY123",
	}
	res, err := svc.ProcessPayment(req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "success", res.Status)
}
