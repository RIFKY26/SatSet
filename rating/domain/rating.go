package domain

type Rating struct {
	ID       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID  string `json:"order_id"`
	DriverID string `json:"driver_id"`
	Score    int    `json:"score"`
	Feedback string `json:"feedback"`
}

type RatingRepository interface {
	SaveRating(rating *Rating) error
}
