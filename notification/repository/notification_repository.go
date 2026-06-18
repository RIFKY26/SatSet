package repository

import (
	"satset2/notification/domain"

	"gorm.io/gorm"
)

type SqlNotificationRepository struct {
	DB *gorm.DB
}

func NewSqlNotificationRepository(db *gorm.DB) domain.NotificationRepository {
	return &SqlNotificationRepository{DB: db}
}

func (r *SqlNotificationRepository) Save(notification *domain.Notification) error {
	return r.DB.Create(notification).Error
}
