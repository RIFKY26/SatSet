package domain

type Promo struct {
	PromoID       string  `json:"promo_id" gorm:"primaryKey"`
	PromoCode     string  `json:"promo_code" gorm:"unique"` // Kode promo tidak boleh kembar
	MinOrderValue float64 `json:"min_order_value"`
	MaxDiscount   float64 `json:"max_discount"`
	DiscountPct   float64 `json:"discount_percent"`
	Quota         int     `json:"quota_remaining"`
}

type PromoRepository interface {
	FindByCode(code string) (*Promo, error)
	UpdateQuota(promo *Promo) error
}
