package tests

import (
	"fmt"
	"satset2/driver-service/domain"
	"satset2/driver-service/service"
	"sync"
	"sync/atomic"
	"testing"
)

// TestAssignOrder_Concurrent_OnlyOneSucceeds spawns 100 goroutines all trying
// to assign an order to the same driver. Exactly one must succeed.
//
// Run with: go test -race ./...
// The -race flag is what proves correctness — a passing count alone is not enough.
func TestAssignOrder_Concurrent_OnlyOneSucceeds(t *testing.T) {
	svc := service.NewDriverService()
	if err := svc.RegisterDriver("d1"); err != nil {
		t.Fatalf("setup: RegisterDriver failed: %v", err)
	}

	const goroutines = 100

	var (
		wg           sync.WaitGroup
		successCount atomic.Int32 // atomic to safely count across goroutines
	)

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			if err := svc.AssignOrder("d1"); err == nil {
				successCount.Add(1)
			}
		}()
	}
	wg.Wait()

	if count := successCount.Load(); count != 1 {
		t.Errorf("expected exactly 1 successful AssignOrder, got %d", count)
	}
}

// TestConcurrent_AssignAndComplete tests the interleaving of AssignOrder and
// CompleteOrder across goroutines. The driver state must always be consistent
// at rest — never simultaneously AVAILABLE and being assigned twice.
func TestConcurrent_AssignAndComplete(t *testing.T) {
	svc := service.NewDriverService()
	if err := svc.RegisterDriver("d1"); err != nil {
		t.Fatalf("setup: RegisterDriver failed: %v", err)
	}

	const cycles = 50
	var wg sync.WaitGroup

	for i := 0; i < cycles; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Assign may fail if another goroutine already assigned — that's fine.
			if err := svc.AssignOrder("d1"); err == nil {
				// Only complete if we were the one who assigned.
				_ = svc.CompleteOrder("d1")
			}
		}()
	}

	wg.Wait()

	// KODE BARU (YANG BENAR)
	driver, err := svc.GetDriver("d1")
	if err != nil {
		t.Fatalf("failed to get driver: %v", err)
	}
	if driver.AvailabilityStatus != domain.AvailabilityAvailable && driver.AvailabilityStatus != domain.AvailabilityOnTrip {
		t.Errorf("driver in invalid state: %s", driver.AvailabilityStatus)
	}
}

// TestConcurrent_MultipleDrivers ensures that concurrent operations on
// different drivers do not interfere with each other.
func TestConcurrent_MultipleDrivers(t *testing.T) {
	svc := service.NewDriverService()

	const numDrivers = 10
	for i := 0; i < numDrivers; i++ {
		id := driverID(i)
		if err := svc.RegisterDriver(id); err != nil {
			t.Fatalf("setup: RegisterDriver(%s) failed: %v", id, err)
		}
	}

	var (
		wg           sync.WaitGroup
		successCount atomic.Int32
	)

	// Each driver gets 10 goroutines racing to assign it.
	for i := 0; i < numDrivers; i++ {
		id := driverID(i)
		for j := 0; j < 10; j++ {
			wg.Add(1)
			go func(driverID string) {
				defer wg.Done()
				if err := svc.AssignOrder(driverID); err == nil {
					successCount.Add(1)
				}
			}(id)
		}
	}

	wg.Wait()

	// Each of the 10 drivers should have exactly 1 successful assignment.
	if count := successCount.Load(); count != numDrivers {
		t.Errorf("expected %d successful assignments (one per driver), got %d",
			numDrivers, count)
	}
}

// TestConcurrent_RegisterAndAssign tests a race between registering a driver
// and immediately assigning an order to it.
func TestConcurrent_RegisterAndAssign(t *testing.T) {
	svc := service.NewDriverService()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		_ = svc.RegisterDriver("d1")
	}()

	go func() {
		defer wg.Done()
		// May fail if register hasn't completed yet — that's valid and expected.
		_ = svc.AssignOrder("d1")
	}()

	wg.Wait()
	// No panic, no data race = success.
}

// driverID is a simple helper to generate driver IDs for table-driven concurrency tests.
func driverID(i int) string {
	return fmt.Sprintf("driver-%d", i)
}
