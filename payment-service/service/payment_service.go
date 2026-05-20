package service

import (
	"fmt"
	"satset2/payment-service/domain"
	"time"
)

type paymentService struct {
	repo domain.PaymentRepository
}

func NewPaymentService(repo domain.PaymentRepository) domain.PaymentService {
	return &paymentService{repo: repo}
}

func (s *paymentService) ProcessPayment(orderID, userID string, amount float64, action string, idmpKey string) (*domain.Transaction, error) {
	// Cek idempotensi
	if existing, _ := s.repo.FindByIdempotencyKey(idmpKey); existing != nil {
		return existing, nil
	}

	status := "pending"
	switch action {
	case "authorize":
		status = "authorized"
	case "capture":
		status = "success"
	}

	tx := &domain.Transaction{
		TransactionID: fmt.Sprintf("PAY-%d", time.Now().UnixNano()),
		OrderID:       orderID,
		UserID:        userID,
		Amount:        amount,
		Status:        status,
		CreatedAt:     time.Now(),
	}

	err := s.repo.CreateTransaction(tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
