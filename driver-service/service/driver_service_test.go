package service

import (
	"satset2/driver-service/domain"
	"testing"
)

// --- RegisterDriver ---

func TestRegisterDriver_Success(t *testing.T) {
	svc := NewDriverService()

	if err := svc.RegisterDriver("d1"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	driver := svc.drivers["d1"]
	if driver.ConnectionStatus != domain.ConnectionOnline {
		t.Errorf("expected ONLINE, got %s", driver.ConnectionStatus)
	}
	if driver.AvailabilityStatus != domain.AvailabilityOnTrip {
		t.Errorf("expected AVAILABLE, got %s", driver.AvailabilityStatus)
	}
}

func TestRegisterDriver_DuplicateID(t *testing.T) {
	svc := NewDriverService()
	_ = svc.RegisterDriver("d1")

	if err := svc.RegisterDriver("d1"); err == nil {
		t.Fatal("expected error for duplicate driver ID, got nil")
	}
}

// --- AssignOrder ---

func TestAssignOrder_Success(t *testing.T) {
	svc := NewDriverService()
	_ = svc.RegisterDriver("d1")

	if err := svc.AssignOrder("d1"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if svc.drivers["d1"].AvailabilityStatus != domain.AvailabilityOnTrip {
		t.Error("expected driver to be ON_TRIP after AssignOrder")
	}
}

func TestAssignOrder_DriverNotFound(t *testing.T) {
	svc := NewDriverService()

	if err := svc.AssignOrder("ghost"); err == nil {
		t.Fatal("expected error for non-existent driver, got nil")
	}
}

func TestAssignOrder_DriverOffline(t *testing.T) {
	svc := NewDriverService()
	_ = svc.RegisterDriver("d1")
	svc.drivers["d1"].ConnectionStatus = domain.ConnectionOffline

	if err := svc.AssignOrder("d1"); err == nil {
		t.Fatal("expected error when assigning order to OFFLINE driver, got nil")
	}
}

func TestAssignOrder_DriverAlreadyOnTrip(t *testing.T) {
	svc := NewDriverService()
	_ = svc.RegisterDriver("d1")
	_ = svc.AssignOrder("d1") // first assignment

	if err := svc.AssignOrder("d1"); err == nil {
		t.Fatal("expected error when assigning order to ON_TRIP driver, got nil")
	}
}

// --- CompleteOrder ---

func TestCompleteOrder_Success(t *testing.T) {
	svc := NewDriverService()
	_ = svc.RegisterDriver("d1")
	_ = svc.AssignOrder("d1")

	if err := svc.CompleteOrder("d1"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if svc.drivers["d1"].AvailabilityStatus != domain.AvailabilityAvailable {
		t.Error("expected driver to be AVAILABLE after CompleteOrder")
	}
}

func TestCompleteOrder_DriverNotFound(t *testing.T) {
	svc := NewDriverService()

	if err := svc.CompleteOrder("ghost"); err == nil {
		t.Fatal("expected error for non-existent driver, got nil")
	}
}

func TestCompleteOrder_DriverNotOnTrip(t *testing.T) {
	svc := NewDriverService()
	_ = svc.RegisterDriver("d1") // starts as AVAILABLE

	if err := svc.CompleteOrder("d1"); err == nil {
		t.Fatal("expected error when completing order for AVAILABLE driver, got nil")
	}
}

// --- Full lifecycle ---

func TestFullLifecycle(t *testing.T) {
	svc := NewDriverService()

	// Register
	if err := svc.RegisterDriver("d1"); err != nil {
		t.Fatalf("RegisterDriver: %v", err)
	}

	// Assign
	if err := svc.AssignOrder("d1"); err != nil {
		t.Fatalf("AssignOrder: %v", err)
	}

	// Complete
	if err := svc.CompleteOrder("d1"); err != nil {
		t.Fatalf("CompleteOrder: %v", err)
	}

	// Driver should be back to AVAILABLE
	if svc.drivers["d1"].AvailabilityStatus != domain.AvailabilityOnTrip {
		t.Error("expected driver back to AVAILABLE after full lifecycle")
	}
}
