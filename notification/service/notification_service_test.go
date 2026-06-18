package service

import (
	"satset2/notification/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockNotificationRepo struct {
	mock.Mock
}

func (m *MockNotificationRepo) Save(notif *domain.Notification) error {
	args := m.Called(notif)
	return args.Error(0)
}

func TestSendNotification_Success(t *testing.T) {
	mockRepo := new(MockNotificationRepo)
	mockRepo.On("Save", mock.Anything).Return(nil)

	svc := NewNotificationService(mockRepo)
	err := svc.SendNotification(1, "Promo", "Diskon 50% pakai SATSET50!")

	assert.NoError(t, err)
}

func TestSendNotification_ErrorValidation(t *testing.T) {
	mockRepo := new(MockNotificationRepo)
	svc := NewNotificationService(mockRepo)

	// Tes gagal jika ID user 0
	err := svc.SendNotification(0, "", "")
	assert.Error(t, err)
}
