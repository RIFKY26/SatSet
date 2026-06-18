package service

import (
	"fmt"
	"satset2/driver-service/domain"
)

type DriverService struct {
	repo domain.DriverRepository
}

func NewDriverService(repo domain.DriverRepository) *DriverService {
	return &DriverService{repo: repo}
}

func (s *DriverService) RegisterDriver(id string) error {
	// Cek apakah driver sudah ada di database
	_, err := s.repo.FindByID(id)
	if err == nil {
		return fmt.Errorf("driver with id %q already registered", id)
	}

	driver := &domain.Driver{
		ID:                 id,
		ConnectionStatus:   domain.ConnectionOnline,
		AvailabilityStatus: domain.AvailabilityAvailable,
	}
	return s.repo.Save(driver)
}

func (s *DriverService) AssignOrder(id string) error {
	driver, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("driver with id %q not found", id)
	}
	if driver.ConnectionStatus != domain.ConnectionOnline {
		return fmt.Errorf("cannot assign order to driver %q: driver is OFFLINE", id)
	}
	if driver.AvailabilityStatus != domain.AvailabilityAvailable {
		return fmt.Errorf("cannot assign order to driver %q: driver is already ON_TRIP", id)
	}

	driver.AvailabilityStatus = domain.AvailabilityOnTrip
	return s.repo.Update(driver)
}

func (s *DriverService) CompleteOrder(id string) error {
	driver, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("driver with id %q not found", id)
	}
	if driver.AvailabilityStatus != domain.AvailabilityOnTrip {
		return fmt.Errorf("cannot complete order for driver %q: driver is not ON_TRIP", id)
	}

	driver.AvailabilityStatus = domain.AvailabilityAvailable
	return s.repo.Update(driver)
}

func (s *DriverService) GetDriver(id string) (*domain.Driver, error) {
	driver, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("driver with id %q not found", id)
	}
	return driver, nil
}