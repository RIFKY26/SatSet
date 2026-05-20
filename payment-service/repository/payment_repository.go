package repository

import "satset2/payment-service/domain"

type paymentRepository struct {
	// Nanti koneksi db (seperti gorm.DB atau *sql.DB) ditaruh di sini
}

func NewPaymentRepository() domain.PaymentRepository {
	return &paymentRepository{}
}

// 1. Memenuhi janji method CreateTransaction
func (r *paymentRepository) CreateTransaction(tx *domain.Transaction) error {
	// Simulasi seolah-olah sukses menyimpan ke database
	return nil
}

// 2. Memenuhi janji method FindByIdempotencyKey
func (r *paymentRepository) FindByIdempotencyKey(key string) (*domain.Transaction, error) {
	// Simulasi seolah-olah belum ada transaksi dengan key ini (tidak duplikat)
	return nil, nil
}
