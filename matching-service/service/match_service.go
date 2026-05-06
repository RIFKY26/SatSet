package service

import (
	"errors"
	"satset2/matching-service/domain"
)

// MatchDriver adalah logika bisnis mencari driver
func MatchDriver(req domain.MatchRequest) (domain.MatchResponse, error) {
	if req.OrderID == "" {
		return domain.MatchResponse{}, errors.New("order ID tidak boleh kosong")
	}

	// Simulasi: selalu berikan driver 123
	return domain.MatchResponse{
		DriverID: 123,
		ETA:      "5 menit",
	}, nil
}
