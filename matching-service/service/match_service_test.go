package service

import (
	"satset2/matching-service/domain"
	"testing"
)

func TestMatchDriver_Success(t *testing.T) {
	req := domain.MatchRequest{OrderID: "ORD-123"}
	result, err := MatchDriver(req)

	if err != nil {
		t.Fatalf("tidak diharapkan error: %v", err)
	}
	if result.DriverID != 123 {
		t.Errorf("driver ID seharusnya 123, dapat %d", result.DriverID)
	}
	if result.ETA == "" {
		t.Error("ETA tidak boleh kosong")
	}
}

func TestMatchDriver_EmptyOrderID(t *testing.T) {
	req := domain.MatchRequest{OrderID: ""}
	_, err := MatchDriver(req)

	if err == nil {
		t.Error("seharusnya error karena OrderID kosong")
	}
}
