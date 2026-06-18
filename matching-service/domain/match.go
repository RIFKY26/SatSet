package domain

type MatchRequest struct {
	OrderID   string  `json:"order_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type MatchResponse struct {
	DriverID string `json:"driver_id"` // Ubah jadi string
	ETA      string `json:"eta"`
}

// Kontrak untuk ngobrol dengan Location Service
type LocationClient interface {
	GetNearestDriver(lat, lng float64) (string, error)
}