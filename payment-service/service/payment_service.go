package service

import (
	"errors"
	"fmt"
	"time"

	"satset2/payment-service/domain"
)

type PaymentRequest struct {
	OrderID        string  `json:"order_id"`
	UserID         int     `json:"user_id"`
	Amount         float64 `json:"amount"`
	PaymentMethod  string  `json:"payment_method"`
	IdempotencyKey string  `json:"idempotency_key"`
}

type PaymentService struct {
	repo domain.PaymentRepository
}

func NewPaymentService(repo domain.PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) ProcessPayment(req PaymentRequest) (*domain.Transaction, error) {
	// 1. Cek Idempotency Key agar saldo tidak terpotong dua kali
	if req.IdempotencyKey != "" {
		existingTx, err := s.repo.FindByIdempotencyKey(req.IdempotencyKey)
		if err == nil && existingTx != nil {
			// Mengembalikan transaksi yang sudah pernah sukses
			return existingTx, nil
		}
	}

	if req.Amount <= 0 {
		return nil, errors.New("jumlah pembayaran tidak valid")
	}

	// 2. Buat Transaksi Baru
	tx := &domain.Transaction{
		TransactionID:  fmt.Sprintf("TXN-%d", time.Now().UnixNano()),
		OrderID:        req.OrderID,
		UserID:         req.UserID,
		Amount:         req.Amount,
		Status:         "success", // Dalam tugas ini, kita asumsikan pembayaran selalu sukses
		PaymentMethod:  req.PaymentMethod,
		IdempotencyKey: req.IdempotencyKey,
		CreatedAt:      time.Now(),
	}

	// 3. Simpan ke Database
	if err := s.repo.CreateTransaction(tx); err != nil {
		return nil, err
	}

	return tx, nil
}
