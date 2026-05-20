package domain

import "time"

// Promo adalah cetak biru untuk data promosi
type Promo struct {
	PromoID         string    `json:"promo_id"`
	PromoCode       string    `json:"promo_code"`
	MinOrderValue   float64   `json:"min_order_value"`
	MaxDiscount     float64   `json:"max_discount"`
	DiscountPercent float64   `json:"discount_percent"`
	QuotaRemaining  int       `json:"quota_remaining"`
	ExpiryDate      time.Time `json:"expiry_date"`
	ServiceType     string    `json:"service_type"` // motor atau mobil
}

// UsageHistory mencatat riwayat penggunaan promo oleh user
type UsageHistory struct {
	ID      string    `json:"id"`
	PromoID string    `json:"promo_id"`
	UserID  string    `json:"user_id"`
	UsedAt  time.Time `json:"used_at"`
}

// PromoRequest untuk input JSON dari handler
type PromoRequest struct {
	PromoCode   string  `json:"promo_code"`
	UserID      string  `json:"user_id"`
	OrderValue  float64 `json:"order_value"`
	ServiceType string  `json:"service_type"`
}

// PromoOutput untuk hasil kalkulasi
type PromoOutput struct {
	IsValid        bool    `json:"is_valid"`
	DiscountAmount float64 `json:"discount_amount"`
	FinalPrice     float64 `json:"final_price"`
	ErrorMessage   string  `json:"error_message"`
}

// PromoRepository adalah kontrak untuk akses data
type PromoRepository interface {
	FindByCode(code string) (*Promo, error)
	GetUserUsageCount(promoID, userID string) (int, error)
	UpdateQuota(promoID string, amount int) error
	RecordUsage(history *UsageHistory) error
}

// PromoService adalah kontrak untuk logika bisnis
type PromoService interface {
	ApplyPromo(code, userID string, orderValue float64, serviceType string) (*PromoOutput, error) // [cite: 70]
}
