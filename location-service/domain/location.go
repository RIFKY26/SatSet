package domain

type DriverLocation struct {
	DriverID  string  `json:"driver_id" gorm:"primaryKey"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type LocationRepository interface {
	UpdateLocation(loc *DriverLocation) error
	GetAllLocations() ([]DriverLocation, error)
}
