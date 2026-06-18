package service

import (
	"math"
	"satset2/location-service/domain"
	"sort"
)

type LocationService struct {
	repo domain.LocationRepository
}

func NewLocationService(repo domain.LocationRepository) *LocationService {
	return &LocationService{repo: repo}
}

func (s *LocationService) UpdateDriverLocation(loc *domain.DriverLocation) error {
	return s.repo.UpdateLocation(loc)
}

// Struct untuk mengembalikan hasil pencarian
type DriverDistance struct {
	DriverID string  `json:"driver_id"`
	Distance float64 `json:"distance"`
}

func (s *LocationService) GetNearestDrivers(lat, lng float64, radius float64) ([]DriverDistance, error) {
	allLocs, err := s.repo.GetAllLocations()
	if err != nil {
		return nil, err
	}

	var results []DriverDistance
	for _, loc := range allLocs {
		// Rumus Pythagoras sederhana untuk simulasi jarak koordinat
		latDiff := loc.Latitude - lat
		lngDiff := loc.Longitude - lng
		distance := math.Sqrt((latDiff*latDiff)+(lngDiff*lngDiff)) * 111 // Konversi kasar ke KM

		if distance <= radius {
			results = append(results, DriverDistance{
				DriverID: loc.DriverID,
				Distance: distance,
			})
		}
	}

	// Urutkan dari yang paling dekat
	sort.Slice(results, func(i, j int) bool {
		return results[i].Distance < results[j].Distance
	})

	return results, nil
}
