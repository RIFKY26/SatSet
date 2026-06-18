package service

import (
	"errors"
	"time"

	"satset2/notification/domain"
)

type NotificationService struct {
	repo domain.NotificationRepository
}

func NewNotificationService(repo domain.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) SendNotification(userID int, title, message string) error {
	if userID <= 0 || title == "" || message == "" {
		return errors.New("user_id, title, dan message tidak boleh kosong")
	}

	notif := &domain.Notification{
		UserID:    userID,
		Title:     title,
		Message:   message,
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	return s.repo.Save(notif)
}
