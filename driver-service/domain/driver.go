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
	ID                 string
	ConnectionStatus   ConnectionStatus
	AvailabilityStatus AvailabilityStatus
}
