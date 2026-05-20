package domain

import "time"

// Transaction mewakili data yang disimpan di DB pusat
type Transaction struct {
	TransactionID string
	OrderID       string
	UserID        string
	Amount        float64
	Status        string
	PaymentMethod string
	CreatedAt     time.Time
}

type PaymentRepository interface {
	CreateTransaction(tx *Transaction) error
	FindByIdempotencyKey(key string) (*Transaction, error)
}

type PaymentService interface {
	ProcessPayment(orderID, userID string, amount float64, action string, idmpKey string) (*Transaction, error)
}
