package domain

type ConnectionStatus string
type AvailabilityStatus string

const (
	ConnectionOnline      ConnectionStatus   = "ONLINE"
	ConnectionOffline     ConnectionStatus   = "OFFLINE"
	AvailabilityAvailable AvailabilityStatus = "AVAILABLE"
	AvailabilityOnTrip    AvailabilityStatus = "ON_TRIP"
)

type Driver struct {
	ID                 string             `json:"id" gorm:"primaryKey"`
	ConnectionStatus   ConnectionStatus   `json:"connection_status"`
	AvailabilityStatus AvailabilityStatus `json:"availability_status"`
}

// Tambahkan kontrak Repository agar strukturnya rapi (Clean Architecture)
type DriverRepository interface {
	Save(driver *Driver) error
	FindByID(id string) (*Driver, error)
	Update(driver *Driver) error
}