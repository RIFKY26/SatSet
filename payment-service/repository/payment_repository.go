package repository

import (
	"satset2/payment-service/domain"

	"gorm.io/gorm"
)

type SqlPaymentRepository struct {
	DB *gorm.DB
}

func NewSqlPaymentRepository(db *gorm.DB) domain.PaymentRepository {
	return &SqlPaymentRepository{DB: db}
}

func (r *SqlPaymentRepository) CreateTransaction(tx *domain.Transaction) error {
	return r.DB.Create(tx).Error
}

func (r *SqlPaymentRepository) FindByIdempotencyKey(key string) (*domain.Transaction, error) {
	var tx domain.Transaction
	// Cari apakah kode pembayaran ini sudah pernah sukses sebelumnya
	err := r.DB.Where("idempotency_key = ?", key).First(&tx).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}
