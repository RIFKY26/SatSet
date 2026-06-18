package domain

import "time"

type Notification struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"` // Tandai apakah notif sudah dibaca
	CreatedAt time.Time `json:"created_at"`
}

type NotificationRepository interface {
	Save(notification *Notification) error
}
