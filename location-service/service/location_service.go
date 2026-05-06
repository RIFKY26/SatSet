package service

import (
	"math"
	"satset2/location-service/domain"
	"sync"
)

type LocationService struct {
	mu           sync.RWMutex
	locations    map[string]domain.Location
	driverClient domain.DriverClient // Menggunakan interface dari domain
}

func NewLocationService(dc domain.DriverClient) *LocationService {
	return &LocationService{
		locations:    make(map[string]domain.Location),
		driverClient: dc,
	}
}

func (s *LocationService) UpdateLocation(driverID string, lat, lng float64, timestamp int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.locations[driverID] = domain.Location{
		DriverID:  driverID,
		Latitude:  lat,
		Longitude: lng,
		Timestamp: timestamp,
	}
}

func (s *LocationService) GetNearbyDrivers(lat, lng, radius float64) []domain.DriverDistance {
	s.mu.RLock()
	defer s.mu.RUnlock()

	results := []domain.DriverDistance{}

	for driverID, loc := range s.locations {
		// Panggil interface untuk cek apakah driver AVAILABLE
		available, err := s.driverClient.IsDriverAvailable(driverID)
		if err != nil || !available {
			continue
		}

		dist := euclideanDistance(lat, lng, loc.Latitude, loc.Longitude)
		if dist <= radius {
			results = append(results, domain.DriverDistance{
				DriverID: driverID,
				Distance: dist,
			})
		}
	}
	return results
}

func euclideanDistance(lat1, lng1, lat2, lng2 float64) float64 {
	dlat := lat1 - lat2
	dlng := lng1 - lng2
	return math.Sqrt(dlat*dlat + dlng*dlng)
}
