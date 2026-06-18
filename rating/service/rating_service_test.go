package service

import (
	"satset2/rating/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRatingRepo struct {
	mock.Mock
}

func (m *MockRatingRepo) SaveRating(rating *domain.Rating) error {
	args := m.Called(rating)
	return args.Error(0)
}

func TestSubmitRating_Success(t *testing.T) {
	mockRepo := new(MockRatingRepo)
	mockRepo.On("SaveRating", mock.Anything).Return(nil)

	svc := NewRatingService(mockRepo)

	rating := &domain.Rating{
		OrderID:  "ORD-1",
		DriverID: "DRV-1",
		Score:    5,
		Feedback: "Mantap sat set!",
	}

	err := svc.SubmitRating(rating)

	assert.NoError(t, err)
}

func TestSubmitRating_InvalidScore(t *testing.T) {
	mockRepo := new(MockRatingRepo)
	svc := NewRatingService(mockRepo)

	rating := &domain.Rating{
		OrderID:  "ORD-1",
		DriverID: "DRV-1",
		Score:    6, // Score tidak valid
	}

	err := svc.SubmitRating(rating)
	assert.Error(t, err)
	assert.Equal(t, "score harus antara 1 sampai 5", err.Error())
}
