package service

import (
	"fmt"
	"satset2/driver-service/domain"
	"sync"
)

type DriverService struct {
	mu      sync.Mutex
	drivers map[string]*domain.Driver
}

func NewDriverService() *DriverService {
	return &DriverService{
		drivers: make(map[string]*domain.Driver),
	}
}

func (s *DriverService) RegisterDriver(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.drivers[id]; exists {
		return fmt.Errorf("driver with id %q already registered", id)
	}

	s.drivers[id] = &domain.Driver{
		ID:                 id,
		ConnectionStatus:   domain.ConnectionOnline,
		AvailabilityStatus: domain.AvailabilityAvailable,
	}
	return nil
}

func (s *DriverService) AssignOrder(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	driver, err := s.findDriver(id)
	if err != nil {
		return err
	}
	if driver.ConnectionStatus != domain.ConnectionOnline {
		return fmt.Errorf("cannot assign order to driver %q: driver is OFFLINE", id)
	}
	if driver.AvailabilityStatus != domain.AvailabilityAvailable {
		return fmt.Errorf("cannot assign order to driver %q: driver is already ON_TRIP", id)
	}

	driver.AvailabilityStatus = domain.AvailabilityOnTrip
	return nil
}

func (s *DriverService) CompleteOrder(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	driver, err := s.findDriver(id)
	if err != nil {
		return err
	}
	if driver.AvailabilityStatus != domain.AvailabilityOnTrip {
		return fmt.Errorf("cannot complete order for driver %q: driver is not ON_TRIP", id)
	}

	driver.AvailabilityStatus = domain.AvailabilityAvailable
	return nil
}

func (s *DriverService) findDriver(id string) (*domain.Driver, error) {
	driver, exists := s.drivers[id]
	if !exists {
		return nil, fmt.Errorf("driver with id %q not found", id)
	}
	return driver, nil
}

func (s *DriverService) GetDriver(id string) (*domain.Driver, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.findDriver(id)
}
