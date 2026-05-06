package service

import (
	"errors"
	"fmt"
	"time"

	"satset2/order-service/domain" // Import struktur data dari domain
)

// CreateOrder memanggil struct dari package domain
func CreateOrder(req domain.OrderRequest, repo domain.OrderRepository) (domain.OrderResponse, error) {
	if req.UserID <= 0 {
		return domain.OrderResponse{}, errors.New("user ID tidak valid")
	}
	if req.Item == "" {
		return domain.OrderResponse{}, errors.New("item tidak boleh kosong")
	}

	orderID := fmt.Sprintf("ORD-%d", time.Now().UnixNano())
	order := domain.OrderResponse{
		OrderID: orderID,
		Status:  "CREATED",
	}

	// Simpan menggunakan repository
	if err := repo.Save(order); err != nil {
		return domain.OrderResponse{}, err
	}

	return order, nil
}
