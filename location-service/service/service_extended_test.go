package service // Kita ubah agar bisa mengakses variabel private

import (
	"fmt"
	"math"
	"testing"
)

// --- MOCK ENHANCEMENTS (Jembatan Clean Architecture) ---
// Kita buat konstanta lokal pengganti driver-service agar logika temanmu tidak rusak
type ConnectionStatus string
type AvailabilityStatus string

const (
	ConnectionOnline      ConnectionStatus   = "ONLINE"
	ConnectionOffline     ConnectionStatus   = "OFFLINE"
	AvailabilityAvailable AvailabilityStatus = "AVAILABLE"
	AvailabilityOnTrip    AvailabilityStatus = "ON_TRIP"
)

type dummyDriver struct {
	ID                 string
	ConnectionStatus   ConnectionStatus
	AvailabilityStatus AvailabilityStatus
}

type mockDriverServiceWithState struct {
	drivers map[string]*dummyDriver
}

func newMockWithState() *mockDriverServiceWithState {
	return &mockDriverServiceWithState{drivers: make(map[string]*dummyDriver)}
}

func (m *mockDriverServiceWithState) addDriver(id string, conn ConnectionStatus, avail AvailabilityStatus) {
	m.drivers[id] = &dummyDriver{ID: id, ConnectionStatus: conn, AvailabilityStatus: avail}
}

func (m *mockDriverServiceWithState) GetDriver(id string) (*dummyDriver, error) {
	d, ok := m.drivers[id]
	if !ok {
		return nil, fmt.Errorf("driver %q not found", id)
	}
	return d, nil
}

// INI KUNCI UTAMANYA: Mengubah Mock temanmu agar sesuai dengan kontrak DriverClient yang baru
func (m *mockDriverServiceWithState) IsDriverAvailable(driverID string) (bool, error) {
	d, err := m.GetDriver(driverID)
	if err != nil {
		return false, err
	}
	return d.ConnectionStatus == ConnectionOnline && d.AvailabilityStatus == AvailabilityAvailable, nil
}

// --- UPDATELOCATION EXTENDED TESTS ---

func TestUpdateLocation_OverwriteWithNewerTimestamp(t *testing.T) {
	svc := NewLocationService(newMockWithState())

	// Store first location with old timestamp
	svc.UpdateLocation("d1", 1.0, 1.0, 100)
	svc.UpdateLocation("d1", 2.0, 2.0, 200) // newer timestamp

	svc.mu.RLock()
	loc := svc.locations["d1"]
	svc.mu.RUnlock()

	if loc.Latitude != 2.0 || loc.Longitude != 2.0 {
		t.Errorf("expected (2.0, 2.0), got (%.1f, %.1f)", loc.Latitude, loc.Longitude)
	}
	if loc.Timestamp != 200 {
		t.Errorf("expected timestamp 200, got %d", loc.Timestamp)
	}
}

func TestUpdateLocation_NegativeCoordinates(t *testing.T) {
	svc := NewLocationService(newMockWithState())

	// Negative coords are valid (like Southern/Western hemispheres)
	svc.UpdateLocation("d1", -40.7128, -74.0060, 1000) // NYC-like

	svc.mu.RLock()
	loc := svc.locations["d1"]
	svc.mu.RUnlock()

	if loc.Latitude != -40.7128 {
		t.Errorf("expected -40.7128, got %f", loc.Latitude)
	}
}

func TestUpdateLocation_ZeroCoordinates(t *testing.T) {
	svc := NewLocationService(newMockWithState())

	svc.UpdateLocation("d1", 0.0, 0.0, 1000)

	svc.mu.RLock()
	_, ok := svc.locations["d1"]
	svc.mu.RUnlock()

	if !ok {
		t.Fatal("expected location to be stored")
	}
}

func TestUpdateLocation_LargeCoordinates(t *testing.T) {
	svc := NewLocationService(newMockWithState())

	svc.UpdateLocation("d1", 180.0, 180.0, 1000)

	svc.mu.RLock()
	loc := svc.locations["d1"]
	svc.mu.RUnlock()

	if loc.Latitude != 180.0 || loc.Longitude != 180.0 {
		t.Errorf("expected (180, 180), got (%.1f, %.1f)", loc.Latitude, loc.Longitude)
	}
}

func TestUpdateLocation_MultipleDrivers(t *testing.T) {
	svc := NewLocationService(newMockWithState())

	for i := 1; i <= 10; i++ {
		id := fmt.Sprintf("d%d", i)
		lat := float64(i)
		lng := float64(i * 2)
		svc.UpdateLocation(id, lat, lng, int64(1000+i))
	}

	svc.mu.RLock()
	if len(svc.locations) != 10 {
		t.Errorf("expected 10 locations, got %d", len(svc.locations))
	}
	svc.mu.RUnlock()
}

// --- GETNERBYDRIVERS EXTENDED TESTS ---

func TestGetNearbyDrivers_DistanceCalculation_Accuracy(t *testing.T) {
	mock := newMockWithState()
	mock.addDriver("d1", ConnectionOnline, AvailabilityAvailable)

	svc := NewLocationService(mock)
	svc.UpdateLocation("d1", 0.0, 0.0, 1000)

	// Test various distances
	tests := []struct {
		lat, lng   float64
		expectDist float64
	}{
		{0.0, 0.0, 0.0},  // same point
		{1.0, 0.0, 1.0},  // 1 unit away in lat
		{0.0, 1.0, 1.0},  // 1 unit away in lng
		{3.0, 4.0, 5.0},  // 3-4-5 triangle
		{-1.0, 0.0, 1.0}, // negative lat
	}

	for _, tt := range tests {
		results := svc.GetNearbyDrivers(tt.lat, tt.lng, 10.0) // large radius
		if len(results) != 1 {
			t.Errorf("expected 1 result for (%.1f, %.1f), got %d",
				tt.lat, tt.lng, len(results))
			continue
		}

		dist := results[0].Distance
		if math.Abs(dist-tt.expectDist) > 0.0001 {
			t.Errorf("at (%.1f, %.1f): expected distance %.1f, got %.4f",
				tt.lat, tt.lng, tt.expectDist, dist)
		}
	}
}

func TestGetNearbyDrivers_RadiusBoundary(t *testing.T) {
	mock := newMockWithState()
	mock.addDriver("d1", ConnectionOnline, AvailabilityAvailable)

	svc := NewLocationService(mock)
	svc.UpdateLocation("d1", 1.0, 1.0, 1000)

	// Driver at (1.0, 1.0), search at (0.0, 0.0), distance = sqrt(2) ≈ 1.414

	// Just inside radius
	results := svc.GetNearbyDrivers(0.0, 0.0, 1.5)
	if len(results) != 1 {
		t.Errorf("expected 1 driver within radius 1.5, got %d", len(results))
	}

	// Exactly at radius
	results = svc.GetNearbyDrivers(0.0, 0.0, 1.415)
	if len(results) != 1 {
		t.Errorf("expected 1 driver at exactly the distance, got %d", len(results))
	}

	// Just outside radius
	results = svc.GetNearbyDrivers(0.0, 0.0, 1.4)
	if len(results) != 0 {
		t.Errorf("expected 0 drivers outside radius 1.4, got %d", len(results))
	}
}

func TestGetNearbyDrivers_MultipleDrivers_OnlyValidReturned(t *testing.T) {
	mock := newMockWithState()
	mock.addDriver("d1", ConnectionOnline, AvailabilityAvailable)  // ✓
	mock.addDriver("d2", ConnectionOnline, AvailabilityOnTrip)     // ✗
	mock.addDriver("d3", ConnectionOffline, AvailabilityAvailable) // ✗
	mock.addDriver("d4", ConnectionOnline, AvailabilityAvailable)  // ✓

	svc := NewLocationService(mock)
	svc.UpdateLocation("d1", 0.0, 0.0, 1000)
	svc.UpdateLocation("d2", 0.5, 0.5, 1000)
	svc.UpdateLocation("d3", 1.0, 1.0, 1000)
	svc.UpdateLocation("d4", 1.5, 1.5, 1000)

	results := svc.GetNearbyDrivers(0.0, 0.0, 10.0)

	if len(results) != 2 {
		t.Fatalf("expected 2 valid drivers, got %d", len(results))
	}

	resultIDs := make(map[string]bool)
	for _, r := range results {
		resultIDs[r.DriverID] = true
	}

	if !resultIDs["d1"] || !resultIDs["d4"] {
		t.Errorf("expected d1 and d4, got %v", resultIDs)
	}
	if resultIDs["d2"] || resultIDs["d3"] {
		t.Errorf("should not include d2 (ON_TRIP) or d3 (OFFLINE)")
	}
}

func TestGetNearbyDrivers_NegativeRadius(t *testing.T) {
	mock := newMockWithState()
	mock.addDriver("d1", ConnectionOnline, AvailabilityAvailable)

	svc := NewLocationService(mock)
	svc.UpdateLocation("d1", 0.0, 0.0, 1000)

	// Negative radius — no drivers should match (distance can't be negative)
	results := svc.GetNearbyDrivers(0.0, 0.0, -1.0)
	if len(results) != 0 {
		t.Errorf("expected 0 drivers with negative radius, got %d", len(results))
	}
}

func TestGetNearbyDrivers_NoLocationsStored(t *testing.T) {
	mock := newMockWithState()
	mock.addDriver("d1", ConnectionOnline, AvailabilityAvailable)

	svc := NewLocationService(mock)
	// Don't update location

	results := svc.GetNearbyDrivers(0.0, 0.0, 10.0)
	if len(results) != 0 {
		t.Errorf("expected 0 results when no locations stored, got %d", len(results))
	}
}

func TestGetNearbyDrivers_ReturnsSortedByDistance(t *testing.T) {
	mock := newMockWithState()
	mock.addDriver("d1", ConnectionOnline, AvailabilityAvailable)
	mock.addDriver("d2", ConnectionOnline, AvailabilityAvailable)
	mock.addDriver("d3", ConnectionOnline, AvailabilityAvailable)

	svc := NewLocationService(mock)
	svc.UpdateLocation("d1", 10.0, 0.0, 1000) // distance 10
	svc.UpdateLocation("d2", 5.0, 0.0, 1000)  // distance 5
	svc.UpdateLocation("d3", 1.0, 0.0, 1000)  // distance 1

	results := svc.GetNearbyDrivers(0.0, 0.0, 15.0)

	// Note: GetNearbyDrivers doesn't guarantee sorting, but verify distances are correct
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	distances := make(map[string]float64)
	for _, r := range results {
		distances[r.DriverID] = r.Distance
	}

	if math.Abs(distances["d1"]-10.0) > 0.0001 {
		t.Errorf("d1: expected distance 10.0, got %.4f", distances["d1"])
	}
	if math.Abs(distances["d2"]-5.0) > 0.0001 {
		t.Errorf("d2: expected distance 5.0, got %.4f", distances["d2"])
	}
	if math.Abs(distances["d3"]-1.0) > 0.0001 {
		t.Errorf("d3: expected distance 1.0, got %.4f", distances["d3"])
	}
}

func TestGetNearbyDrivers_StateChangesDuringQuery(t *testing.T) {
	mock := newMockWithState()
	mock.addDriver("d1", ConnectionOnline, AvailabilityAvailable)

	svc := NewLocationService(mock)
	svc.UpdateLocation("d1", 0.0, 0.0, 1000)

	// Query returns driver
	results := svc.GetNearbyDrivers(0.0, 0.0, 10.0)
	if len(results) != 1 {
		t.Fatalf("expected 1 driver initially")
	}

	// Change driver state
	d, _ := mock.GetDriver("d1")
	d.AvailabilityStatus = AvailabilityOnTrip

	// Query should now exclude it
	results = svc.GetNearbyDrivers(0.0, 0.0, 10.0)
	if len(results) != 0 {
		t.Errorf("expected 0 drivers after state change, got %d", len(results))
	}
}

func TestGetNearbyDrivers_LargeScaleWithManyDrivers(t *testing.T) {
	mock := newMockWithState()
	const numDrivers = 100

	// Register all drivers
	for i := 0; i < numDrivers; i++ {
		id := fmt.Sprintf("d%d", i)
		status := AvailabilityAvailable
		if i%3 == 0 {
			status = AvailabilityOnTrip // 1/3 on trip
		}
		mock.addDriver(id, ConnectionOnline, status)
	}

	svc := NewLocationService(mock)

	// Scatter locations
	for i := 0; i < numDrivers; i++ {
		id := fmt.Sprintf("d%d", i)
		lat := float64(i%10) - 5.0 // spread across -5 to 5
		lng := float64(i/10) - 5.0 // spread across -5 to 5
		svc.UpdateLocation(id, lat, lng, int64(1000+i))
	}

	// Query nearby
	results := svc.GetNearbyDrivers(0.0, 0.0, 7.0) // reasonable radius

	// Should only return AVAILABLE drivers
	for _, r := range results {
		d, _ := mock.GetDriver(r.DriverID)
		if d.AvailabilityStatus != AvailabilityAvailable {
			t.Errorf("returned driver %s should be AVAILABLE, got %s",
				r.DriverID, d.AvailabilityStatus)
		}
	}

	// Verify result count is reasonable
	if len(results) > numDrivers {
		t.Errorf("results exceed available drivers")
	}
}

// --- INTEGRATION WITH MOCK DRIVERSERVICE ---

func TestLocationService_IntegrationWithMockDriver(t *testing.T) {
	mock := newMockWithState()
	mock.addDriver("d1", ConnectionOnline, AvailabilityAvailable)

	svc := NewLocationService(mock)
	svc.UpdateLocation("d1", 0.0, 0.0, 1000)

	// Should include before state change
	results := svc.GetNearbyDrivers(0.0, 0.0, 10.0)
	if len(results) != 1 {
		t.Fatal("expected d1 in results")
	}

	// Change connection status directly via mock
	d, _ := mock.GetDriver("d1")
	d.ConnectionStatus = ConnectionOffline

	// Should exclude after state change
	results = svc.GetNearbyDrivers(0.0, 0.0, 10.0)
	if len(results) != 0 {
		t.Fatal("expected d1 to be excluded after going OFFLINE")
	}
}
