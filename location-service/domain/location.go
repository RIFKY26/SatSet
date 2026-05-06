package domain

type Location struct {
	DriverID  string
	Latitude  float64
	Longitude float64
	Timestamp int64 // Unix timestamp in seconds
}

type DriverDistance struct {
	DriverID string  `json:"driver_id"`
	Distance float64 `json:"distance"`
}

// DriverClient adalah kontrak/aturan untuk siapa saja yang mau mengecek status driver
type DriverClient interface {
	IsDriverAvailable(driverID string) (bool, error)
}
