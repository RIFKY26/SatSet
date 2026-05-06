package tests

import (
	"fmt"
	"satset2/driver-service/domain"
	"satset2/driver-service/service"
	"strings"
	"testing"
)

// --- TEST HELPERS ---

// verifyDriverState checks that a driver is in the expected state.
// Used to validate driver state after operations.
func verifyDriverState(t *testing.T, svc *service.DriverService, driverID string, expectedConn domain.ConnectionStatus, expectedAvail domain.AvailabilityStatus) {
	d, err := svc.GetDriver(driverID)
	if err != nil {
		t.Fatalf("failed to get driver %q: %v", driverID, err)
	}
	if d.ConnectionStatus != expectedConn {
		t.Errorf("driver %q: expected connection %s, got %s", driverID, expectedConn, d.ConnectionStatus)
	}
	if d.AvailabilityStatus != expectedAvail {
		t.Errorf("driver %q: expected availability %s, got %s", driverID, expectedAvail, d.AvailabilityStatus)
	}
}

// --- REGISTERDRIVER EXTENDED TESTS ---

func TestRegisterDriver_MultipleValidDrivers(t *testing.T) {
	svc := service.NewDriverService()

	for i := 1; i <= 5; i++ {
		id := fmt.Sprintf("d%d", i)
		if err := svc.RegisterDriver(id); err != nil {
			t.Fatalf("RegisterDriver(%s) failed: %v", id, err)
		}
	}

	// All should be ONLINE + AVAILABLE
	for i := 1; i <= 5; i++ {
		id := fmt.Sprintf("d%d", i)
		verifyDriverState(t, svc, id, domain.ConnectionOnline, domain.AvailabilityAvailable)
	}
}

func TestRegisterDriver_EmptyID(t *testing.T) {
	svc := service.NewDriverService()
	// Empty ID is allowed by the service — it's just a key in the map
	// The application layer would validate this, not the service
	if err := svc.RegisterDriver(""); err != nil {
		t.Fatalf("RegisterDriver with empty ID should succeed: %v", err)
	}
}

func TestRegisterDriver_LongID(t *testing.T) {
	svc := service.NewDriverService()
	longID := strings.Repeat("a", 1000)

	if err := svc.RegisterDriver(longID); err != nil {
		t.Fatalf("RegisterDriver with long ID failed: %v", err)
	}
	verifyDriverState(t, svc, longID, domain.ConnectionOnline, domain.AvailabilityAvailable)
}

func TestRegisterDriver_SpecialCharactersInID(t *testing.T) {
	svc := service.NewDriverService()
	specialIDs := []string{
		"d-1",
		"d_2",
		"d.3",
		"d@4",
		"d#5",
		"d$6",
		"d 7",  // space
		"d\n8", // newline
	}

	for _, id := range specialIDs {
		if err := svc.RegisterDriver(id); err != nil {
			t.Fatalf("RegisterDriver(%q) failed: %v", id, err)
		}
	}
}

// --- ASSIGNORDER EXTENDED TESTS ---

func TestAssignOrder_SequentialAssignments(t *testing.T) {
	svc := service.NewDriverService()
	_ = svc.RegisterDriver("d1")

	// First assignment succeeds
	if err := svc.AssignOrder("d1"); err != nil {
		t.Fatalf("first AssignOrder failed: %v", err)
	}
	verifyDriverState(t, svc, "d1", domain.ConnectionOnline, domain.AvailabilityOnTrip)

	// Second assignment fails (already ON_TRIP)
	err := svc.AssignOrder("d1")
	if err == nil {
		t.Fatal("second AssignOrder should fail, got nil")
	}

	// Error message is descriptive
	if !strings.Contains(err.Error(), "ON_TRIP") {
		t.Errorf("expected 'ON_TRIP' in error message, got: %v", err)
	}
}

func TestAssignOrder_AfterCompleteOrder_CanAssignAgain(t *testing.T) {
	svc := service.NewDriverService()
	_ = svc.RegisterDriver("d1")

	// Assign
	_ = svc.AssignOrder("d1")
	verifyDriverState(t, svc, "d1", domain.ConnectionOnline, domain.AvailabilityOnTrip)

	// Complete
	_ = svc.CompleteOrder("d1")
	verifyDriverState(t, svc, "d1", domain.ConnectionOnline, domain.AvailabilityAvailable)

	// Can assign again
	if err := svc.AssignOrder("d1"); err != nil {
		t.Fatalf("reassignment after completion failed: %v", err)
	}
	verifyDriverState(t, svc, "d1", domain.ConnectionOnline, domain.AvailabilityOnTrip)
}

func TestAssignOrder_MultipleDrivers_Independent(t *testing.T) {
	svc := service.NewDriverService()
	_ = svc.RegisterDriver("d1")
	_ = svc.RegisterDriver("d2")

	// Assign d1
	_ = svc.AssignOrder("d1")

	// d1 should be ON_TRIP, d2 should still be AVAILABLE
	verifyDriverState(t, svc, "d1", domain.ConnectionOnline, domain.AvailabilityOnTrip)
	verifyDriverState(t, svc, "d2", domain.ConnectionOnline, domain.AvailabilityAvailable)

	// Can still assign d2
	if err := svc.AssignOrder("d2"); err != nil {
		t.Fatalf("AssignOrder(d2) failed: %v", err)
	}
	verifyDriverState(t, svc, "d2", domain.ConnectionOnline, domain.AvailabilityOnTrip)
}

func TestAssignOrder_NonExistentDriver(t *testing.T) {
	svc := service.NewDriverService()

	err := svc.AssignOrder("ghost")
	if err == nil {
		t.Fatal("expected error for non-existent driver")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' in error, got: %v", err)
	}
}

func TestAssignOrder_OfflineDriver(t *testing.T) {
	svc := service.NewDriverService()
	_ = svc.RegisterDriver("d1")

	// Manually set OFFLINE
	d, _ := svc.GetDriver("d1")
	d.ConnectionStatus = domain.ConnectionOffline

	err := svc.AssignOrder("d1")
	if err == nil {
		t.Fatal("expected error when assigning to OFFLINE driver")
	}

	if !strings.Contains(err.Error(), "OFFLINE") {
		t.Errorf("expected 'OFFLINE' in error, got: %v", err)
	}
}

// --- COMPLETEORDER EXTENDED TESTS ---

func TestCompleteOrder_SequentialOperations(t *testing.T) {
	svc := service.NewDriverService()
	_ = svc.RegisterDriver("d1")

	// Complete before assign should fail
	if err := svc.CompleteOrder("d1"); err == nil {
		t.Fatal("expected error when completing before assignment")
	}

	// Assign
	_ = svc.AssignOrder("d1")
	verifyDriverState(t, svc, "d1", domain.ConnectionOnline, domain.AvailabilityOnTrip)

	// Complete succeeds
	if err := svc.CompleteOrder("d1"); err != nil {
		t.Fatalf("CompleteOrder failed: %v", err)
	}
	verifyDriverState(t, svc, "d1", domain.ConnectionOnline, domain.AvailabilityAvailable)

	// Second complete fails
	if err := svc.CompleteOrder("d1"); err == nil {
		t.Fatal("expected error on second CompleteOrder")
	}
}

func TestCompleteOrder_AfterOfflineTransition(t *testing.T) {
	svc := service.NewDriverService()
	_ = svc.RegisterDriver("d1")
	_ = svc.AssignOrder("d1")

	// Go offline while ON_TRIP
	d, _ := svc.GetDriver("d1")
	d.ConnectionStatus = domain.ConnectionOffline

	// Should still be able to complete (no ONLINE check in CompleteOrder)
	// This allows drivers to recover from disconnects
	if err := svc.CompleteOrder("d1"); err != nil {
		t.Fatalf("CompleteOrder should work even while OFFLINE: %v", err)
	}
	verifyDriverState(t, svc, "d1", domain.ConnectionOffline, domain.AvailabilityAvailable)
}

func TestCompleteOrder_NonExistentDriver(t *testing.T) {
	svc := service.NewDriverService()

	if err := svc.CompleteOrder("ghost"); err == nil {
		t.Fatal("expected error for non-existent driver")
	}
}

// --- GETDRIVER TESTS ---

func TestGetDriver_ReturnsCorrectState(t *testing.T) {
	svc := service.NewDriverService()
	_ = svc.RegisterDriver("d1")

	d, err := svc.GetDriver("d1")
	if err != nil {
		t.Fatalf("GetDriver failed: %v", err)
	}

	if d.ID != "d1" {
		t.Errorf("expected ID d1, got %s", d.ID)
	}
	if d.ConnectionStatus != domain.ConnectionOnline {
		t.Errorf("expected ONLINE, got %s", d.ConnectionStatus)
	}
	if d.AvailabilityStatus != domain.AvailabilityAvailable {
		t.Errorf("expected AVAILABLE, got %s", d.AvailabilityStatus)
	}
}

func TestGetDriver_NonExistent(t *testing.T) {
	svc := service.NewDriverService()

	_, err := svc.GetDriver("ghost")
	if err == nil {
		t.Fatal("expected error for non-existent driver")
	}
}

func TestGetDriver_ReflectsStateChanges(t *testing.T) {
	svc := service.NewDriverService()
	_ = svc.RegisterDriver("d1")

	// Initial state
	d1, _ := svc.GetDriver("d1")
	if d1.AvailabilityStatus != domain.AvailabilityAvailable {
		t.Error("initial state should be AVAILABLE")
	}

	// After assign
	_ = svc.AssignOrder("d1")
	d2, _ := svc.GetDriver("d1")
	if d2.AvailabilityStatus != domain.AvailabilityOnTrip {
		t.Error("state after assignment should be ON_TRIP")
	}

	// After complete
	_ = svc.CompleteOrder("d1")
	d3, _ := svc.GetDriver("d1")
	if d3.AvailabilityStatus != domain.AvailabilityAvailable {
		t.Error("state after completion should be AVAILABLE")
	}
}

// --- ERROR MESSAGE VALIDATION ---

func TestErrorMessages_AreDescriptive(t *testing.T) {
	svc := service.NewDriverService()

	tests := []struct {
		name       string
		operation  func() error
		wantSubstr string
	}{
		{
			name: "RegisterDriver duplicate",
			operation: func() error {
				_ = svc.RegisterDriver("d1")
				return svc.RegisterDriver("d1")
			},
			wantSubstr: "already registered",
		},
		{
			name: "AssignOrder not found",
			operation: func() error {
				return svc.AssignOrder("ghost")
			},
			wantSubstr: "not found",
		},
		{
			name: "CompleteOrder not on trip",
			operation: func() error {
				_ = svc.RegisterDriver("d2")
				return svc.CompleteOrder("d2")
			},
			wantSubstr: "not ON_TRIP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.operation()
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), tt.wantSubstr) {
				t.Errorf("error %q missing %q", err.Error(), tt.wantSubstr)
			}
		})
	}
}

// --- FULL LIFECYCLE TESTS ---

func TestFullLifecycle_MultipleRounds(t *testing.T) {
	svc := service.NewDriverService()
	_ = svc.RegisterDriver("d1")

	for round := 0; round < 5; round++ {
		// Assign
		if err := svc.AssignOrder("d1"); err != nil {
			t.Fatalf("round %d: AssignOrder failed: %v", round, err)
		}

		// Verify ON_TRIP
		d, _ := svc.GetDriver("d1")
		if d.AvailabilityStatus != domain.AvailabilityOnTrip {
			t.Errorf("round %d: expected ON_TRIP", round)
		}

		// Complete
		if err := svc.CompleteOrder("d1"); err != nil {
			t.Fatalf("round %d: CompleteOrder failed: %v", round, err)
		}

		// Verify AVAILABLE
		d, _ = svc.GetDriver("d1")
		if d.AvailabilityStatus != domain.AvailabilityAvailable {
			t.Errorf("round %d: expected AVAILABLE after completion", round)
		}
	}
}

func TestFullLifecycle_ManyDrivers(t *testing.T) {
	svc := service.NewDriverService()
	const numDrivers = 20

	// Register all
	for i := 0; i < numDrivers; i++ {
		id := fmt.Sprintf("d%d", i)
		if err := svc.RegisterDriver(id); err != nil {
			t.Fatalf("RegisterDriver(%s) failed: %v", id, err)
		}
	}

	// Assign all
	for i := 0; i < numDrivers; i++ {
		id := fmt.Sprintf("d%d", i)
		if err := svc.AssignOrder(id); err != nil {
			t.Fatalf("AssignOrder(%s) failed: %v", id, err)
		}
	}

	// Verify all ON_TRIP
	for i := 0; i < numDrivers; i++ {
		id := fmt.Sprintf("d%d", i)
		d, _ := svc.GetDriver(id)
		if d.AvailabilityStatus != domain.AvailabilityOnTrip {
			t.Errorf("driver %s should be ON_TRIP", id)
		}
	}

	// Complete all
	for i := 0; i < numDrivers; i++ {
		id := fmt.Sprintf("d%d", i)
		if err := svc.CompleteOrder(id); err != nil {
			t.Fatalf("CompleteOrder(%s) failed: %v", id, err)
		}
	}

	// Verify all AVAILABLE
	for i := 0; i < numDrivers; i++ {
		id := fmt.Sprintf("d%d", i)
		d, _ := svc.GetDriver(id)
		if d.AvailabilityStatus != domain.AvailabilityAvailable {
			t.Errorf("driver %s should be AVAILABLE", id)
		}
	}
}
