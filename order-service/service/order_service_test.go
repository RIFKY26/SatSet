package service

import (
	"testing"

	"satset2/order-service/domain" // Import struktur data dari domain
)

// MockOrderRepository adalah implementasi palsu untuk testing
type MockOrderRepository struct {
	SaveFunc func(order domain.OrderResponse) error
}

func (m *MockOrderRepository) Save(order domain.OrderResponse) error {
	return m.SaveFunc(order)
}

// TestCreateOrder_Success menguji pembuatan order berhasil
func TestCreateOrder_Success(t *testing.T) {
	// Setup mock
	mockRepo := &MockOrderRepository{
		SaveFunc: func(order domain.OrderResponse) error {
			// Pastikan status adalah CREATED
			if order.Status != "CREATED" {
				t.Errorf("status tidak sesuai, dapat %s", order.Status)
			}
			return nil
		},
	}

	req := domain.OrderRequest{UserID: 1, Item: "Makanan"}
	order, err := CreateOrder(req, mockRepo)

	if err != nil {
		t.Fatalf("seharusnya tidak error, dapat: %v", err)
	}
	if order.OrderID == "" {
		t.Error("OrderID kosong")
	}
	if order.Status != "CREATED" {
		t.Errorf("status seharusnya CREATED, dapat %s", order.Status)
	}
}

// TestCreateOrder_InvalidUserID menguji validasi user ID
func TestCreateOrder_InvalidUserID(t *testing.T) {
	mockRepo := &MockOrderRepository{} // tidak dipanggil
	req := domain.OrderRequest{UserID: 0, Item: "Makanan"}

	_, err := CreateOrder(req, mockRepo)

	if err == nil {
		t.Error("seharusnya error karena user ID 0")
	}
}

// TestCreateOrder_EmptyItem menguji item kosong
func TestCreateOrder_EmptyItem(t *testing.T) {
	mockRepo := &MockOrderRepository{}
	req := domain.OrderRequest{UserID: 1, Item: ""}

	_, err := CreateOrder(req, mockRepo)

	if err == nil {
		t.Error("seharusnya error karena item kosong")
	}
}
