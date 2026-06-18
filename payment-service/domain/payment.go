package domain

import "time"

// Ini adalah tabel transactions yang akan dibuat di database
type Transaction struct {
	TransactionID  string    `json:"transaction_id" gorm:"primaryKey"`
	OrderID        string    `json:"order_id"`
	UserID         int       `json:"user_id"`
	Amount         float64   `json:"amount"`
	Status         string    `json:"status"` // pending, success, failed
	PaymentMethod  string    `json:"payment_method"`
	IdempotencyKey string    `json:"idempotency_key" gorm:"unique"` // Mencegah bayar double
	CreatedAt      time.Time `json:"created_at"`
}

type PaymentRepository interface {
	CreateTransaction(tx *Transaction) error
	FindByIdempotencyKey(key string) (*Transaction, error)
}
