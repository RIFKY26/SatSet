package service

import (
	"satset2/location-service/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLocationRepo struct {
	mock.Mock
}

func (m *MockLocationRepo) UpdateLocation(loc *domain.DriverLocation) error {
	args := m.Called(loc)
	return args.Error(0)
}

func (m *MockLocationRepo) GetAllLocations() ([]domain.DriverLocation, error) {
	args := m.Called()
	return args.Get(0).([]domain.DriverLocation), args.Error(1)
}

func TestUpdateDriverLocation(t *testing.T) {
	mockRepo := new(MockLocationRepo)
	mockRepo.On("UpdateLocation", mock.Anything).Return(nil)

	svc := NewLocationService(mockRepo)
	loc := &domain.DriverLocation{DriverID: "DRV-1", Latitude: -6.2, Longitude: 106.8}
	err := svc.UpdateDriverLocation(loc)

	assert.NoError(t, err)
}

func TestGetNearestDrivers(t *testing.T) {
	mockRepo := new(MockLocationRepo)
	mockLocs := []domain.DriverLocation{
		{DriverID: "DRV-123", Latitude: -6.2001, Longitude: 106.8167},
	}
	mockRepo.On("GetAllLocations").Return(mockLocs, nil)

	svc := NewLocationService(mockRepo)
	// Cari dalam radius 5 KM
	res, err := svc.GetNearestDrivers(-6.2000, 106.8166, 5.0)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "DRV-123", res[0].DriverID)
}
