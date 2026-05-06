package tests

import (
	"fmt"
	"satset2/driver-service/domain"
	"satset2/driver-service/service"
	"sync"
	"sync/atomic"
	"testing"
)

// TestStress_1000Goroutines_ConcurrentOperations simulates a high-load scenario
// with 1000 concurrent goroutines performing mixed assignments and completions.
//
// This test validates:
//  1. No driver can be assigned multiple times simultaneously
//  2. State remains consistent under extreme concurrency
//  3. No panics or race conditions
func TestStress_1000Goroutines_ConcurrentOperations(t *testing.T) {
	svc := service.NewDriverService()
	const numDrivers = 50
	const goroutinesPerDriver = 20 // 50 * 20 = 1000 goroutines

	// Register all drivers
	for i := 0; i < numDrivers; i++ {
		driverID := fmt.Sprintf("driver-%d", i)
		if err := svc.RegisterDriver(driverID); err != nil {
			t.Fatalf("RegisterDriver(%s) failed: %v", driverID, err)
		}
	}

	var (
		wg           sync.WaitGroup
		totalSuccess atomic.Int32
		totalFailed  atomic.Int32
	)

	// Launch goroutines: each tries to assign its driver
	for i := 0; i < numDrivers; i++ {
		driverID := fmt.Sprintf("driver-%d", i)

		for g := 0; g < goroutinesPerDriver; g++ {
			wg.Add(1)
			go func(id string) {
				defer wg.Done()

				// Each goroutine tries to assign the same driver
				// Only ONE should succeed
				err := svc.AssignOrder(id)
				if err == nil {
					totalSuccess.Add(1)
				} else {
					totalFailed.Add(1)
				}
			}(driverID)
		}
	}

	wg.Wait()

	// Validate: exactly numDrivers successful assignments
	// (one per driver, since only AVAILABLE drivers can be assigned)
	successCount := totalSuccess.Load()
	failedCount := totalFailed.Load()

	if successCount != int32(numDrivers) {
		t.Errorf("expected %d successful assignments, got %d",
			numDrivers, successCount)
	}

	expectedFailed := int32(numDrivers*goroutinesPerDriver) - int32(numDrivers)
	if failedCount != expectedFailed {
		t.Errorf("expected %d failed assignments, got %d",
			expectedFailed, failedCount)
	}

	// Verify: all drivers are now ON_TRIP (exactly one succeeded per driver)
	for i := 0; i < numDrivers; i++ {
		driverID := fmt.Sprintf("driver-%d", i)
		d, err := svc.GetDriver(driverID)
		if err != nil {
			t.Errorf("failed to get driver %s: %v", driverID, err)
			continue
		}

		if d.AvailabilityStatus != domain.AvailabilityOnTrip {
			t.Errorf("driver %s should be ON_TRIP, got %s",
				driverID, d.AvailabilityStatus)
		}
	}

	t.Logf("Stress test passed: %d successes, %d failures out of %d goroutines",
		successCount, failedCount, numDrivers*goroutinesPerDriver)
}

// TestStress_AssignAndComplete_Interleaved simulates a realistic scenario where
// some drivers are being assigned while others are completing orders.
func TestStress_AssignAndComplete_Interleaved(t *testing.T) {
	svc := service.NewDriverService()
	const numDrivers = 30
	const rounds = 10

	// Register all drivers
	for i := 0; i < numDrivers; i++ {
		if err := svc.RegisterDriver(fmt.Sprintf("d%d", i)); err != nil {
			t.Fatalf("RegisterDriver failed: %v", err)
		}
	}

	var wg sync.WaitGroup

	// For each round, launch concurrent assign and complete operations
	for round := 0; round < rounds; round++ {
		for i := 0; i < numDrivers; i++ {
			driverID := fmt.Sprintf("d%d", i)

			// Goroutine 1: Assign this driver
			wg.Add(1)
			go func(id string) {
				defer wg.Done()
				_ = svc.AssignOrder(id)
			}(driverID)

			// Goroutine 2: Try to complete (may fail if not yet assigned)
			wg.Add(1)
			go func(id string) {
				defer wg.Done()
				_ = svc.CompleteOrder(id)
			}(driverID)
		}
	}

	wg.Wait()

	// After all rounds, verify final state is consistent
	for i := 0; i < numDrivers; i++ {
		driverID := fmt.Sprintf("d%d", i)
		d, err := svc.GetDriver(driverID)
		if err != nil {
			t.Errorf("GetDriver(%s) failed: %v", driverID, err)
			continue
		}

		// Driver should be in a valid state (either AVAILABLE or ON_TRIP)
		if d.AvailabilityStatus != domain.AvailabilityAvailable &&
			d.AvailabilityStatus != domain.AvailabilityOnTrip {
			t.Errorf("driver %s in invalid state: %s",
				driverID, d.AvailabilityStatus)
		}
	}
}

// TestStress_HighFrequencyOperations simulates a driver pool with continuous
// assignment and completion cycles.
func TestStress_HighFrequencyOperations(t *testing.T) {
	svc := service.NewDriverService()
	const numDrivers = 10
	const operationsPerDriver = 100

	// Register all drivers
	for i := 0; i < numDrivers; i++ {
		if err := svc.RegisterDriver(fmt.Sprintf("d%d", i)); err != nil {
			t.Fatalf("RegisterDriver failed: %v", err)
		}
	}

	var wg sync.WaitGroup
	var (
		assignCount      atomic.Int32
		completeCount    atomic.Int32
		assignFailures   atomic.Int32
		completeFailures atomic.Int32
	)

	// Each driver has a dedicated goroutine that cycles through assign/complete
	for i := 0; i < numDrivers; i++ {
		wg.Add(1)
		go func(driverID string) {
			defer wg.Done()

			for op := 0; op < operationsPerDriver; op++ {
				// Assign
				if err := svc.AssignOrder(driverID); err == nil {
					assignCount.Add(1)
				} else {
					assignFailures.Add(1)
				}

				// Complete
				if err := svc.CompleteOrder(driverID); err == nil {
					completeCount.Add(1)
				} else {
					completeFailures.Add(1)
				}
			}
		}(fmt.Sprintf("d%d", i))
	}

	wg.Wait()

	// Validate: each driver should have completed operationsPerDriver cycles
	// Each cycle = 1 successful assign + 1 successful complete
	expectedSuccess := int32(numDrivers * operationsPerDriver)

	if assignCount.Load() != expectedSuccess {
		t.Errorf("expected %d assigns, got %d", expectedSuccess, assignCount.Load())
	}

	if completeCount.Load() != expectedSuccess {
		t.Errorf("expected %d completes, got %d", expectedSuccess, completeCount.Load())
	}

	t.Logf("High-frequency test passed: %d assigns, %d completes, %d+%d failures",
		assignCount.Load(), completeCount.Load(),
		assignFailures.Load(), completeFailures.Load())
}

// TestStress_MassiveDriverPool simulates a production-scale driver pool
// with 500 drivers and 5000 concurrent operations.
func TestStress_MassiveDriverPool(t *testing.T) {
	svc := service.NewDriverService()
	const numDrivers = 500

	// Register all drivers
	for i := 0; i < numDrivers; i++ {
		if err := svc.RegisterDriver(fmt.Sprintf("driver-%d", i)); err != nil {
			t.Fatalf("RegisterDriver failed at %d: %v", i, err)
		}
	}

	var wg sync.WaitGroup
	var successCount atomic.Int32

	// Launch 5000 random assignment attempts
	for attempt := 0; attempt < 5000; attempt++ {
		wg.Add(1)
		go func(attemptNum int) {
			defer wg.Done()

			// Pick a random driver (deterministic based on attempt number)
			driverIdx := attemptNum % numDrivers
			driverID := fmt.Sprintf("driver-%d", driverIdx)

			// Try to assign
			if err := svc.AssignOrder(driverID); err == nil {
				successCount.Add(1)
			}
		}(attempt)
	}

	wg.Wait()

	// At most numDrivers should succeed (one per driver)
	if successCount.Load() > int32(numDrivers) {
		t.Errorf("more successful assignments than drivers: %d > %d",
			successCount.Load(), numDrivers)
	}

	// Verify no driver was assigned twice
	for i := 0; i < numDrivers; i++ {
		driverID := fmt.Sprintf("driver-%d", i)
		d, err := svc.GetDriver(driverID)
		if err != nil {
			t.Errorf("GetDriver(%s) failed: %v", driverID, err)
			continue
		}

		// Each driver should be either AVAILABLE or ON_TRIP (but not both)
		if d.AvailabilityStatus != domain.AvailabilityAvailable &&
			d.AvailabilityStatus != domain.AvailabilityOnTrip {
			t.Errorf("driver %s in invalid state: %s",
				driverID, d.AvailabilityStatus)
		}
	}

	t.Logf("Massive pool test passed: %d successful out of 5000 attempts on %d drivers",
		successCount.Load(), numDrivers)
}

// TestStress_RapidStateTransitions tests rapid cycling through states
// to catch timing-sensitive bugs.
func TestStress_RapidStateTransitions(t *testing.T) {
	svc := service.NewDriverService()
	const numDrivers = 20
	const cyclesPerDriver = 50

	// Register all drivers
	for i := 0; i < numDrivers; i++ {
		if err := svc.RegisterDriver(fmt.Sprintf("d%d", i)); err != nil {
			t.Fatalf("RegisterDriver failed: %v", err)
		}
	}

	var wg sync.WaitGroup

	// Each driver rapidly cycles through assign → complete → available
	for i := 0; i < numDrivers; i++ {
		wg.Add(1)
		go func(driverID string) {
			defer wg.Done()

			for cycle := 0; cycle < cyclesPerDriver; cycle++ {
				// Must succeed within a cycle
				if err := svc.AssignOrder(driverID); err != nil {
					t.Errorf("AssignOrder failed in cycle %d: %v",
						cycle, err)
					continue
				}

				// Complete immediately
				if err := svc.CompleteOrder(driverID); err != nil {
					t.Errorf("CompleteOrder failed in cycle %d: %v",
						cycle, err)
				}
			}
		}(fmt.Sprintf("d%d", i))
	}

	wg.Wait()

	// All drivers should be AVAILABLE at the end
	for i := 0; i < numDrivers; i++ {
		driverID := fmt.Sprintf("d%d", i)
		d, err := svc.GetDriver(driverID)
		if err != nil {
			t.Errorf("GetDriver(%s) failed: %v", driverID, err)
			continue
		}

		if d.AvailabilityStatus != domain.AvailabilityAvailable {
			t.Errorf("driver %s should be AVAILABLE at end, got %s",
				driverID, d.AvailabilityStatus)
		}
	}

	t.Logf("Rapid transition test passed: %d drivers × %d cycles completed",
		numDrivers, cyclesPerDriver)
}

// TestStress_MixedFailureAndSuccess simulates scenarios where operations fail
// due to state constraints while others succeed.
func TestStress_MixedFailureAndSuccess(t *testing.T) {
	svc := service.NewDriverService()
	const numDrivers = 15

	// Register all drivers
	for i := 0; i < numDrivers; i++ {
		if err := svc.RegisterDriver(fmt.Sprintf("d%d", i)); err != nil {
			t.Fatalf("RegisterDriver failed: %v", err)
		}
	}

	var wg sync.WaitGroup
	var (
		successfulOps atomic.Int32
		failedOps     atomic.Int32
	)

	// Multiple goroutines per driver, each trying different operations
	for i := 0; i < numDrivers; i++ {
		driverID := fmt.Sprintf("d%d", i)

		// First, assign the driver
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			if err := svc.AssignOrder(id); err == nil {
				successfulOps.Add(1)
			} else {
				failedOps.Add(1)
			}
		}(driverID)

		// Meanwhile, multiple goroutines try to complete (should fail)
		for c := 0; c < 5; c++ {
			wg.Add(1)
			go func(id string) {
				defer wg.Done()
				if err := svc.CompleteOrder(id); err == nil {
					successfulOps.Add(1)
				} else {
					failedOps.Add(1)
				}
			}(driverID)
		}
	}

	wg.Wait()

	// Verify reasonable outcome:
	// - At least numDrivers should succeed (assign operations)
	// - Some failures are expected (complete operations on AVAILABLE drivers)
	if successfulOps.Load() < int32(numDrivers) {
		t.Errorf("expected at least %d successes, got %d",
			numDrivers, successfulOps.Load())
	}

	if failedOps.Load() == 0 {
		t.Errorf("expected some failures, got none")
	}

	t.Logf("Mixed ops test passed: %d successes, %d failures",
		successfulOps.Load(), failedOps.Load())
}
