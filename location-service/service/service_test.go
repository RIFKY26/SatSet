package service

import (
	"fmt"
	"testing"
)

// --- MOCK ENHANCEMENTS (Jembatan) ---
type mockDriverService struct {
	drivers map[string]*dummyDriver
}

func newMock() *mockDriverService {
	return &mockDriverService{drivers: make(map[string]*dummyDriver)}
}

func (m *mockDriverService) addDriver(id string, conn ConnectionStatus, avail AvailabilityStatus) {
	m.drivers[id] = &dummyDriver{ID: id, ConnectionStatus: conn, AvailabilityStatus: avail}
}

func (m *mockDriverService) GetDriver(id string) (*dummyDriver, error) {
	d, ok := m.drivers[id]
	if !ok {
		return nil, fmt.Errorf("driver %q not found", id)
	}
	return d, nil
}

func (m *mockDriverService) IsDriverAvailable(driverID string) (bool, error) {
	d, err := m.GetDriver(driverID)
	if err != nil {
		return false, err
	}
	return d.ConnectionStatus == ConnectionOnline && d.AvailabilityStatus == AvailabilityAvailable, nil
}

// --- UpdateLocation ---

func TestUpdateLocation_StoresLocation(t *testing.T) {
	svc := NewLocationService(newMock())
	svc.UpdateLocation("d1", 1.0, 2.0, 1000)

	svc.mu.RLock()
	loc, ok := svc.locations["d1"]
	svc.mu.RUnlock()

	if !ok {
		t.Fatal("expected location to be stored")
	}
	if loc.Latitude != 1.0 || loc.Longitude != 2.0 {
		t.Errorf("unexpected location: %+v", loc)
	}
}

func TestUpdateLocation_OverwritesPreviousLocation(t *testing.T) {
	svc := NewLocationService(newMock())
	svc.UpdateLocation("d1", 1.0, 2.0, 1000)
	svc.UpdateLocation("d1", 9.0, 8.0, 2000) // overwrite

	svc.mu.RLock()
	loc := svc.locations["d1"]
	svc.mu.RUnlock()

	if loc.Latitude != 9.0 || loc.Longitude != 8.0 {
		t.Errorf("expected overwritten location, got: %+v", loc)
	}
	if loc.Timestamp != 2000 {
		t.Errorf("expected timestamp 2000, got %d", loc.Timestamp)
	}
}

// --- GetNearbyDrivers ---

func TestGetNearbyDrivers_ReturnsAvailableDriversInRadius(t *testing.T) {
	mock := newMock()
	mock.addDriver("d1", ConnectionOnline, AvailabilityAvailable)

	svc := NewLocationService(mock)
	svc.UpdateLocation("d1", 0.0, 0.0, 1000)

	results := svc.GetNearbyDrivers(0.0, 0.0, 1.0)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].DriverID != "d1" {
		t.Errorf("unexpected driver: %s", results[0].DriverID)
	}
}

func TestGetNearbyDrivers_ExcludesOfflineDriver(t *testing.T) {
	mock := newMock()
	mock.addDriver("d1", ConnectionOffline, AvailabilityAvailable)

	svc := NewLocationService(mock)
	svc.UpdateLocation("d1", 0.0, 0.0, 1000)

	results := svc.GetNearbyDrivers(0.0, 0.0, 10.0)

	if len(results) != 0 {
		t.Errorf("expected 0 results for OFFLINE driver, got %d", len(results))
	}
}

func TestGetNearbyDrivers_ExcludesOnTripDriver(t *testing.T) {
	mock := newMock()
	mock.addDriver("d1", ConnectionOnline, AvailabilityOnTrip)

	svc := NewLocationService(mock)
	svc.UpdateLocation("d1", 0.0, 0.0, 1000)

	results := svc.GetNearbyDrivers(0.0, 0.0, 10.0)

	if len(results) != 0 {
		t.Errorf("expected 0 results for ON_TRIP driver, got %d", len(results))
	}
}

func TestGetNearbyDrivers_ExcludesDriverOutsideRadius(t *testing.T) {
	mock := newMock()
	mock.addDriver("d1", ConnectionOnline, AvailabilityAvailable)

	svc := NewLocationService(mock)
	svc.UpdateLocation("d1", 100.0, 100.0, 1000) // far away

	results := svc.GetNearbyDrivers(0.0, 0.0, 1.0)

	if len(results) != 0 {
		t.Errorf("expected 0 results for out-of-radius driver, got %d", len(results))
	}
}

func TestGetNearbyDrivers_ExcludesDriverWithNoLocation(t *testing.T) {
	mock := newMock()
	mock.addDriver("d1", ConnectionOnline, AvailabilityAvailable)
	// No UpdateLocation call — d1 has no location

	svc := NewLocationService(mock)
	results := svc.GetNearbyDrivers(0.0, 0.0, 10.0)

	if len(results) != 0 {
		t.Errorf("expected 0 results for driver with no location, got %d", len(results))
	}
}

func TestGetNearbyDrivers_ExcludesDriverNotInDriverService(t *testing.T) {
	mock := newMock()
	// d1 has a location but is NOT registered in DriverService

	svc := NewLocationService(mock)
	svc.UpdateLocation("d1", 0.0, 0.0, 1000)

	results := svc.GetNearbyDrivers(0.0, 0.0, 10.0)

	if len(results) != 0 {
		t.Errorf("expected 0 results for unknown driver, got %d", len(results))
	}
}

func TestGetNearbyDrivers_RadiusZero_OnlyExactMatch(t *testing.T) {
	mock := newMock()
	mock.addDriver("d1", ConnectionOnline, AvailabilityAvailable)

	svc := NewLocationService(mock)
	svc.UpdateLocation("d1", 1.0, 1.0, 1000)

	// Exact position — distance is 0, should be included
	results := svc.GetNearbyDrivers(1.0, 1.0, 0.0)
	if len(results) != 1 {
		t.Errorf("expected 1 result for exact position with radius=0, got %d", len(results))
	}

	// Slightly off — should be excluded
	results = svc.GetNearbyDrivers(1.1, 1.1, 0.0)
	if len(results) != 0 {
		t.Errorf("expected 0 results when radius=0 and not exact match, got %d", len(results))
	}
}
func TestGetNearbyDrivers_ReturnsEmptySliceNotNil(t *testing.T) {
	svc := NewLocationService(newMock())
	results := svc.GetNearbyDrivers(0.0, 0.0, 10.0)

	if results == nil {
		t.Error("expected empty slice, got nil")
	}
}
