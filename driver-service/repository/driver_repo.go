package repository

import (
	"satset2/driver-service/domain"
	"gorm.io/gorm"
)

type SqlDriverRepository struct {
	DB *gorm.DB
}

func NewSqlDriverRepository(db *gorm.DB) domain.DriverRepository {
	return &SqlDriverRepository{DB: db}
}

func (r *SqlDriverRepository) Save(driver *domain.Driver) error {
	return r.DB.Create(driver).Error
}

func (r *SqlDriverRepository) FindByID(id string) (*domain.Driver, error) {
	var driver domain.Driver
	err := r.DB.Where("id = ?", id).First(&driver).Error
	return &driver, err
}

func (r *SqlDriverRepository) Update(driver *domain.Driver) error {
	return r.DB.Save(driver).Error
}