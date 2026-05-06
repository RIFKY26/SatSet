package domain

type MatchRequest struct {
	OrderID string `json:"order_id"`
}

type MatchResponse struct {
	DriverID int    `json:"driver_id"`
	ETA      string `json:"eta"`
}
