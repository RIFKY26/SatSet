package service

import (
	"satset2/notification/domain"
)

type NotificationService struct {
	Provider domain.NotificationProvider
}

// Tambahkan fungsi New (Constructor)
func NewNotificationService(p domain.NotificationProvider) *NotificationService {
	return &NotificationService{Provider: p}
}

// Fungsi ini sudah diimplementasikan untuk memanggil Provider
func (s *NotificationService) SendNotification(userID, message string) error {
	// Di Unit Test akan memanggil Mock, di Functional memanggil layanan asli
	return s.Provider.Send(userID, message)
}
