package service

import (
	"errors"
	"fmt"
	"time"

	"satset2/order-service/domain"
)

func CreateOrder(req domain.OrderRequest, repo domain.OrderRepository) (domain.OrderResponse, error) {
	if req.UserID <= 0 {
		return domain.OrderResponse{}, errors.New("user ID tidak valid")
	}
	if req.Item == "" {
		return domain.OrderResponse{}, errors.New("item tidak boleh kosong")
	}

	// 1. Buat ID Pesanan
	orderID := fmt.Sprintf("ORD-%d", time.Now().UnixNano())
	
	// 2. Siapkan data yang akan dimasukkan ke Database
	order := domain.Order{
		OrderID: orderID,
		UserID:  req.UserID,
		Item:    req.Item,
		Status:  "CREATED",
	}

	// 3. Simpan ke Database asli
	if err := repo.Save(&order); err != nil {
		return domain.OrderResponse{}, err
	}

	// 4. Kembalikan respon ke User
	return domain.OrderResponse{
		OrderID: order.OrderID,
		Status:  order.Status,
	}, nil
}