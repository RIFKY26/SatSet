package domain

// Kontrak untuk nembak Database (Postgres/Mongo)
type RatingRepository interface {
	SaveRating(orderID string, driverID string, score int) error
}
