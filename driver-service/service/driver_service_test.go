package service

import (
	"satset2/driver-service/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDriverRepository struct {
	mock.Mock
}

func (m *MockDriverRepository) Save(driver *domain.Driver) error {
	args := m.Called(driver)
	return args.Error(0)
}

func (m *MockDriverRepository) FindByID(id string) (*domain.Driver, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Driver), args.Error(1)
}

func (m *MockDriverRepository) Update(driver *domain.Driver) error {
	args := m.Called(driver)
	return args.Error(0)
}

func TestRegisterDriver_Success(t *testing.T) {
	mockRepo := new(MockDriverRepository)
	// Pura-puranya driver belum ada di database
	mockRepo.On("FindByID", "DRV-1").Return(nil, assert.AnError)
	mockRepo.On("Save", mock.Anything).Return(nil)

	svc := NewDriverService(mockRepo)
	err := svc.RegisterDriver("DRV-1")

	assert.NoError(t, err)
}

func TestAssignOrder_Success(t *testing.T) {
	mockRepo := new(MockDriverRepository)
	// Pura-puranya driver ada dan sedang nganggur
	driver := &domain.Driver{ID: "DRV-1", ConnectionStatus: domain.ConnectionOnline, AvailabilityStatus: domain.AvailabilityAvailable}
	mockRepo.On("FindByID", "DRV-1").Return(driver, nil)
	mockRepo.On("Update", mock.Anything).Return(nil)

	svc := NewDriverService(mockRepo)
	err := svc.AssignOrder("DRV-1")

	assert.NoError(t, err)
}
