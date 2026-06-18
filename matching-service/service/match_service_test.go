package service

import (
	"satset2/matching-service/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Membuat mock pengganti HTTP client
type MockLocationClient struct {
	mock.Mock
}

func (m *MockLocationClient) GetNearestDriver(lat, lng float64) (string, error) {
	args := m.Called(lat, lng)
	return args.String(0), args.Error(1)
}

func TestMatchDriverSuccess(t *testing.T) {
	mockClient := new(MockLocationClient)
	// Pura-puranya location-service menjawab DRV-123 ada di dekat user
	mockClient.On("GetNearestDriver", -6.200000, 106.816666).Return("DRV-123", nil)

	svc := NewMatchService(mockClient)
	req := domain.MatchRequest{
		OrderID:   "ORD-001",
		Latitude:  -6.200000,
		Longitude: 106.816666,
	}

	res, err := svc.MatchDriver(req)

	assert.NoError(t, err)
	assert.Equal(t, "DRV-123", res.DriverID)
}
