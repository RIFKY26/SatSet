package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProvider struct {
	mock.Mock
}

func (m *MockProvider) Send(userID string, message string) error {
	args := m.Called(userID, message)
	return args.Error(0)
}

func TestSendNotification_UnitTest(t *testing.T) {
	mockProvider := new(MockProvider)
	svc := NotificationService{Provider: mockProvider}

	mockProvider.On("Send", "USER-123", "Ada promo nih!").Return(nil)
	err := svc.SendNotification("USER-123", "Ada promo nih!")

	assert.NoError(t, err, "Seharusnya tidak ada error saat mengirim notifikasi")
}
