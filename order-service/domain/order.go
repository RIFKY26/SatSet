package domain

// Ini adalah bentuk data yang akan disimpan di Database PostgreSQL
type Order struct {
	OrderID string `json:"order_id" gorm:"primaryKey"`
	UserID  int    `json:"user_id"`
	Item    string `json:"item"`
	Status  string `json:"status"` // CREATED, ASSIGNED, COMPLETED
}

// OrderRequest adalah struktur permintaan dari Postman/Frontend
type OrderRequest struct {
	UserID int    `json:"user_id"`
	Item   string `json:"item"`
}

// OrderResponse adalah balasan setelah order dibuat
type OrderResponse struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

// OrderRepository adalah kontrak untuk menyimpan order
type OrderRepository interface {
	Save(order *Order) error // Kita ubah menerima pointer ke struct Order yang baru
}