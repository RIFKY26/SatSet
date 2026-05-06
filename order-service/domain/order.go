package domain

// OrderRequest adalah struktur permintaan
type OrderRequest struct {
	UserID int    `json:"user_id"`
	Item   string `json:"item"`
}

// OrderResponse adalah respons setelah order dibuat
type OrderResponse struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

// OrderRepository adalah interface untuk menyimpan order
type OrderRepository interface {
	Save(order OrderResponse) error
}
