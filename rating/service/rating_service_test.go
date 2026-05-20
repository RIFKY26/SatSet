package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRatingDB struct {
	mock.Mock
}

func (m *MockRatingDB) SaveRating(orderID string, driverID string, score int) error {
	args := m.Called(orderID, driverID, score)
	return args.Error(0)
}

func TestSubmitRating_UnitTest(t *testing.T) {
	mockDB := new(MockRatingDB)
	svc := RatingService{Repo: mockDB}

	mockDB.On("SaveRating", "ORD-123", "DRV-999", 5).Return(nil)
	err := svc.SubmitRating("ORD-123", "DRV-999", 5)

	assert.NoError(t, err, "Seharusnya tidak ada error saat submit rating")
}
