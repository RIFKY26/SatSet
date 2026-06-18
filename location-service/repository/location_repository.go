package repository

import (
	"satset2/location-service/domain"

	"gorm.io/gorm"
)

type SqlLocationRepository struct {
	DB *gorm.DB
}

func NewSqlLocationRepository(db *gorm.DB) domain.LocationRepository {
	return &SqlLocationRepository{DB: db}
}

func (r *SqlLocationRepository) UpdateLocation(loc *domain.DriverLocation) error {
	// Gunakan Save agar jika DriverID sudah ada, datanya di-update, bukan duplikat
	return r.DB.Save(loc).Error
}

func (r *SqlLocationRepository) GetAllLocations() ([]domain.DriverLocation, error) {
	var locations []domain.DriverLocation
	err := r.DB.Find(&locations).Error
	return locations, err
}
