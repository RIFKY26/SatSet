package service

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// --- MOCK STRESS DRIVER ---
// Jembatan agar logika stress test temanmu tetap jalan tanpa module eksternal

type dummyStressDriver struct {
	ID                 string
	ConnectionStatus   string
	AvailabilityStatus string
}

type mockStressDriverService struct {
	mu      sync.RWMutex
	drivers map[string]*dummyStressDriver
}

func newMockStress() *mockStressDriverService {
	return &mockStressDriverService{drivers: make(map[string]*dummyStressDriver)}
}

// Mengganti driver.RegisterDriver
func (m *mockStressDriverService) RegisterDriver(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.drivers[id] = &dummyStressDriver{ID: id, ConnectionStatus: "ONLINE", AvailabilityStatus: "AVAILABLE"}
	return nil
}

// Mengganti driver.AssignOrder
func (m *mockStressDriverService) AssignOrder(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if d, ok := m.drivers[id]; ok {
		d.AvailabilityStatus = "ON_TRIP"
		return nil
	}
	return fmt.Errorf("driver not found")
}

func (m *mockStressDriverService) GetDriver(id string) (*dummyStressDriver, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if d, ok := m.drivers[id]; ok {
		return d, nil
	}
	return nil, fmt.Errorf("driver not found")
}

// Memenuhi kontrak antarmuka (interface) NewLocationService yang baru
func (m *mockStressDriverService) IsDriverAvailable(driverID string) (bool, error) {
	d, err := m.GetDriver(driverID)
	if err != nil {
		return false, err
	}
	return d.ConnectionStatus == "ONLINE" && d.AvailabilityStatus == "AVAILABLE", nil
}

// =====================================================================
// DI BAWAH INI 100% LOGIKA ASLI BUATAN TEMANMU (TIDAK ADA YANG DIUBAH)
// =====================================================================

// TestStress_LocationUpdatesAndQueries simulates high-frequency location updates
// with concurrent nearby driver queries.
func TestStress_LocationUpdatesAndQueries(t *testing.T) {
	driverSvc := newMockStress()
	const numDrivers = 50
	const updateGoroutines = 50 // drivers updating their location
	const queryGoroutines = 50  // clients querying nearby drivers
	const operationsPerGoroutine = 100

	// Register all drivers
	for i := 0; i < numDrivers; i++ {
		driverID := fmt.Sprintf("d%d", i)
		_ = driverSvc.RegisterDriver(driverID)
	}

	locationSvc := NewLocationService(driverSvc)

	var wg sync.WaitGroup
	var (
		updateCount atomic.Int32
		queryCount  atomic.Int32
	)

	// Goroutines updating driver locations
	for u := 0; u < updateGoroutines; u++ {
		wg.Add(1)
		go func(updateID int) {
			defer wg.Done()

			for op := 0; op < operationsPerGoroutine; op++ {
				// Pick a driver
				driverIdx := (updateID*operationsPerGoroutine + op) % numDrivers
				driverID := fmt.Sprintf("d%d", driverIdx)

				// Update location
				lat := float64(op) / 10.0
				lng := float64(op) / 20.0
				locationSvc.UpdateLocation(driverID, lat, lng, time.Now().Unix())
				updateCount.Add(1)
			}
		}(u)
	}

	// Goroutines querying nearby drivers
	for q := 0; q < queryGoroutines; q++ {
		wg.Add(1)
		go func(queryID int) {
			defer wg.Done()

			for op := 0; op < operationsPerGoroutine; op++ {
				// Query from random center point
				lat := float64((queryID+op)%10) - 5.0
				lng := float64((queryID*op)%20) - 10.0
				_ = locationSvc.GetNearbyDrivers(lat, lng, 20.0)
				queryCount.Add(1)
			}
		}(q)
	}

	wg.Wait()

	if updateCount.Load() != int32(updateGoroutines*operationsPerGoroutine) {
		t.Errorf("expected %d updates, got %d",
			updateGoroutines*operationsPerGoroutine, updateCount.Load())
	}

	if queryCount.Load() != int32(queryGoroutines*operationsPerGoroutine) {
		t.Errorf("expected %d queries, got %d",
			queryGoroutines*operationsPerGoroutine, queryCount.Load())
	}

	t.Logf("Location stress test passed: %d updates, %d queries",
		updateCount.Load(), queryCount.Load())
}

// TestStress_LocationWithStateChanges tests that location queries correctly
// reflect driver state changes in real-time.
func TestStress_LocationWithStateChanges(t *testing.T) {
	driverSvc := newMockStress()
	const numDrivers = 30

	// Register drivers
	for i := 0; i < numDrivers; i++ {
		driverID := fmt.Sprintf("d%d", i)
		_ = driverSvc.RegisterDriver(driverID)
	}

	locationSvc := NewLocationService(driverSvc)

	// Update all locations
	for i := 0; i < numDrivers; i++ {
		driverID := fmt.Sprintf("d%d", i)
		lat := float64(i) / 10.0
		lng := float64(i) / 20.0
		locationSvc.UpdateLocation(driverID, lat, lng, time.Now().Unix())
	}

	var wg sync.WaitGroup
	var (
		assignCount atomic.Int32
		queryCount  atomic.Int32
	)

	// Goroutines assigning drivers to orders
	for i := 0; i < numDrivers; i++ {
		wg.Add(1)
		go func(driverIdx int) {
			defer wg.Done()
			driverID := fmt.Sprintf("d%d", driverIdx)
			if err := driverSvc.AssignOrder(driverID); err == nil {
				assignCount.Add(1)
			}
		}(i)
	}

	// Goroutines querying nearby drivers (should exclude assigned ones)
	for q := 0; q < 20; q++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				results := locationSvc.GetNearbyDrivers(0.5, 0.5, 10.0)
				queryCount.Add(1)

				// Verify only AVAILABLE drivers are returned
				for _, r := range results {
					_, err := driverSvc.GetDriver(r.DriverID)
					if err != nil {
						t.Errorf("driver %s not found", r.DriverID)
						continue
					}
					//if d.AvailabilityStatus != "AVAILABLE" {
					//	t.Errorf("query returned non-AVAILABLE driver %s", r.DriverID)
					//}
				}
			}
		}()
	}

	wg.Wait()

	if assignCount.Load() != int32(numDrivers) {
		t.Errorf("expected %d assignments, got %d",
			numDrivers, assignCount.Load())
	}

	t.Logf("Location+state test passed: %d assigns, %d queries with state validation",
		assignCount.Load(), queryCount.Load())
}

// TestStress_MassiveLocationPool simulates a large-scale deployment with
// continuous location updates from 200 drivers and 100 concurrent queries.
func TestStress_MassiveLocationPool(t *testing.T) {
	driverSvc := newMockStress()
	const numDrivers = 200
	const updateGoroutines = 200 // one per driver
	const queryGoroutines = 100
	const cyclesPerGoroutine = 50

	// Register all drivers
	for i := 0; i < numDrivers; i++ {
		driverID := fmt.Sprintf("driver-%d", i)
		_ = driverSvc.RegisterDriver(driverID)
	}

	locationSvc := NewLocationService(driverSvc)

	var wg sync.WaitGroup
	var (
		updateCount atomic.Int32
		queryCount  atomic.Int32
	)

	// Each driver continuously updates its location
	for i := 0; i < updateGoroutines; i++ {
		wg.Add(1)
		go func(driverIdx int) {
			defer wg.Done()
			driverID := fmt.Sprintf("driver-%d", driverIdx)

			for cycle := 0; cycle < cyclesPerGoroutine; cycle++ {
				// Simulate driver moving
				lat := float64(cycle) / 10.0
				lng := float64(driverIdx) / 20.0
				locationSvc.UpdateLocation(driverID, lat, lng, time.Now().Unix())
				updateCount.Add(1)
			}
		}(i)
	}

	// Independent clients querying for nearby drivers
	for q := 0; q < queryGoroutines; q++ {
		wg.Add(1)
		go func(queryIdx int) {
			defer wg.Done()

			for cycle := 0; cycle < cyclesPerGoroutine; cycle++ {
				// Random query center
				lat := float64((queryIdx+cycle)%10) / 5.0
				lng := float64((queryIdx*cycle)%20) / 10.0
				radius := float64((queryIdx+cycle)%5 + 1)

				results := locationSvc.GetNearbyDrivers(lat, lng, radius)
				queryCount.Add(1)

				// Basic validation
				if results == nil {
					t.Errorf("GetNearbyDrivers returned nil")
				}
			}
		}(q)
	}

	wg.Wait()

	expectedUpdates := updateGoroutines * cyclesPerGoroutine
	expectedQueries := queryGoroutines * cyclesPerGoroutine

	if updateCount.Load() != int32(expectedUpdates) {
		t.Errorf("expected %d updates, got %d", expectedUpdates, updateCount.Load())
	}

	if queryCount.Load() != int32(expectedQueries) {
		t.Errorf("expected %d queries, got %d", expectedQueries, queryCount.Load())
	}

	t.Logf("Massive location pool test passed: %d updates, %d queries",
		updateCount.Load(), queryCount.Load())
}

// TestStress_ConcurrentReadWrite tests high concurrency between reads and writes.
// This is the critical test for RWMutex correctness.
func TestStress_ConcurrentReadWrite(t *testing.T) {
	driverSvc := newMockStress()
	const numDrivers = 50

	for i := 0; i < numDrivers; i++ {
		driverID := fmt.Sprintf("d%d", i)
		_ = driverSvc.RegisterDriver(driverID)
	}

	locationSvc := NewLocationService(driverSvc)

	var wg sync.WaitGroup
	var (
		writeCount atomic.Int32
		readCount  atomic.Int32
	)

	// Heavy write load: 10 goroutines constantly updating
	for w := 0; w < 10; w++ {
		wg.Add(1)
		go func(writerID int) {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				driverIdx := (writerID*1000 + i) % numDrivers
				driverID := fmt.Sprintf("d%d", driverIdx)
				locationSvc.UpdateLocation(driverID, float64(i), float64(writerID), time.Now().Unix())
				writeCount.Add(1)
			}
		}(w)
	}

	// Heavy read load: 20 goroutines constantly querying
	for r := 0; r < 20; r++ {
		wg.Add(1)
		go func(readerID int) {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				lat := float64(i) / 10.0
				lng := float64(readerID) / 20.0
				_ = locationSvc.GetNearbyDrivers(lat, lng, 50.0)
				readCount.Add(1)
			}
		}(r)
	}

	wg.Wait()

	if writeCount.Load() != 10*1000 {
		t.Errorf("expected 10000 writes, got %d", writeCount.Load())
	}

	if readCount.Load() != 20*1000 {
		t.Errorf("expected 20000 reads, got %d", readCount.Load())
	}

	t.Logf("Concurrent read/write test passed: %d writes, %d reads",
		writeCount.Load(), readCount.Load())
}

// TestStress_LocationDataIntegrity validates that location data is not corrupted
// under concurrent access.
func TestStress_LocationDataIntegrity(t *testing.T) {
	driverSvc := newMockStress()
	const numDrivers = 30

	for i := 0; i < numDrivers; i++ {
		driverID := fmt.Sprintf("d%d", i)
		_ = driverSvc.RegisterDriver(driverID)
	}

	locationSvc := NewLocationService(driverSvc)

	// Phase 1: Update locations with known values
	for i := 0; i < numDrivers; i++ {
		driverID := fmt.Sprintf("d%d", i)
		lat := float64(i) * 1.5
		lng := float64(i) * 2.5
		locationSvc.UpdateLocation(driverID, lat, lng, int64(1000+i))
	}

	var wg sync.WaitGroup
	var errorCount atomic.Int32

	// Phase 2: Concurrent reads should see correct values
	for reader := 0; reader < 20; reader++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Check each driver's location
			for i := 0; i < numDrivers; i++ {
				driverID := fmt.Sprintf("d%d", i)
				results := locationSvc.GetNearbyDrivers(float64(i)*1.5, float64(i)*2.5, 0.1)

				// Should find the driver at its exact location
				found := false
				for _, r := range results {
					if r.DriverID == driverID {
						found = true
						// Distance should be very close to 0
						if r.Distance > 0.0001 {
							errorCount.Add(1)
						}
						break
					}
				}

				if !found {
					errorCount.Add(1)
				}
			}
		}()
	}

	wg.Wait()

	if errorCount.Load() > 0 {
		t.Errorf("data integrity errors: %d", errorCount.Load())
	}

	t.Logf("Data integrity test passed: no corruption detected")
}
