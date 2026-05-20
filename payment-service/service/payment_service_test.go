package service

import (
	"satset2/payment-service/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateTransaction(t *domain.Transaction) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockRepository) FindByIdempotencyKey(key string) (*domain.Transaction, error) {
	args := m.Called(key)
	var tx *domain.Transaction
	if args.Get(0) != nil {
		tx = args.Get(0).(*domain.Transaction)
	}
	return tx, args.Error(1)
}

func TestProcessPaymentSuccess(t *testing.T) {
	mockRepo := new(MockRepository)

	// Tambahkan baris ini untuk mengajari Mock menjawab pencarian Idempotency Key
	// Kita suruh dia mereturn "nil, nil" (artinya tidak ada duplikat, tidak ada error)
	mockRepo.On("FindByIdempotencyKey", "key-1").Return(nil, nil)

	mockRepo.On("CreateTransaction", mock.Anything).Return(nil)

	svc := NewPaymentService(mockRepo)
	res, err := svc.ProcessPayment("ORD-123", "USER-1", 50000, "authorize", "key-1")

	assert.NoError(t, err)
	assert.Equal(t, "authorized", res.Status)
	mockRepo.AssertExpectations(t)
}
