package repository

import (
	"errors"
	"satset2/notification/domain"
)

type SqlNotificationRepository struct {
	// Nanti diisi: db *sql.DB (koneksi database)
}

func NewNotificationRepository() domain.NotificationRepository {
	return &SqlNotificationRepository{}
}

func (r *SqlNotificationRepository) SaveLog(userID string, status string) error {
	return errors.New("TDD: repository SaveLog belum diimplementasikan")
}
