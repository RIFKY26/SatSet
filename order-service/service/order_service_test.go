package service

import (
	"satset2/order-service/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderRepository struct {
	mock.Mock
}

// Gunakan parameter struct database asli (*domain.Order)
func (m *MockOrderRepository) Save(order *domain.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func TestCreateOrderSuccess(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	mockRepo.On("Save", mock.Anything).Return(nil)

	req := domain.OrderRequest{
		UserID: 1,
		Item:   "Nasi Goreng",
	}

	res, err := CreateOrder(req, mockRepo)

	assert.NoError(t, err)
	assert.NotEmpty(t, res.OrderID)
	assert.Equal(t, "CREATED", res.Status)
}
