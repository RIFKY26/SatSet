package domain

// Kontrak untuk nembak API luar (Firebase/SMS)
type NotificationProvider interface {
	Send(userID string, message string) error
}

// Kontrak untuk nembak Database (Postgres/Mongo)
type NotificationRepository interface {
	SaveLog(userID string, status string) error
}
