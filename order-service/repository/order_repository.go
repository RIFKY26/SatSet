package repository

import (
	"satset2/order-service/domain"
	"gorm.io/gorm"
)

type SqlOrderRepository struct {
	DB *gorm.DB
}

func NewSqlOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &SqlOrderRepository{DB: db}
}

// Fungsi Save sekarang menembak langsung ke PostgreSQL
func (r *SqlOrderRepository) Save(order *domain.Order) error {
	return r.DB.Create(order).Error
}