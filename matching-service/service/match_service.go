package service

import (
	"errors"
	"fmt"
	"satset2/matching-service/domain"
)

type MatchService struct {
	locationClient domain.LocationClient
}

func NewMatchService(lc domain.LocationClient) *MatchService {
	return &MatchService{locationClient: lc}
}

func (s *MatchService) MatchDriver(req domain.MatchRequest) (domain.MatchResponse, error) {
	if req.OrderID == "" {
		return domain.MatchResponse{}, errors.New("order ID tidak boleh kosong")
	}

	// Cari driver ke location-service secara sungguhan!
	driverID, err := s.locationClient.GetNearestDriver(req.Latitude, req.Longitude)
	if err != nil {
		return domain.MatchResponse{}, fmt.Errorf("maaf: %v", err)
	}

	return domain.MatchResponse{
		DriverID: driverID,
		ETA:      "5 menit", // Untuk tugas ini ETA bisa dibiarkan statis dulu
	}, nil
}