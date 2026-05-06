package repository

import (
	"satset2/order-service/domain" // Import struktur data dari domain
)

// DefaultOrderRepository adalah implementasi konkrit
type DefaultOrderRepository struct{}

// Fungsi Save sekarang menerima domain.OrderResponse
func (r DefaultOrderRepository) Save(order domain.OrderResponse) error {
	// Simulasi simpan ke database
	return nil
}
