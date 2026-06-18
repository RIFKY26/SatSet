package repository

import (
	"user_service/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetByID(id string) (*domain.User, error)
	// Tambahkan CreateUser untuk testing nanti
	CreateUser(user *domain.User) error 
}

type SqlUserRepository struct {
	DB *gorm.DB
}

func NewSqlUserRepository(db *gorm.DB) UserRepository {
	return &SqlUserRepository{DB: db}
}

func (r *SqlUserRepository) GetByID(id string) (*domain.User, error) {
	var user domain.User
	result := r.DB.Where("id = ?", id).First(&user) 
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *SqlUserRepository) CreateUser(user *domain.User) error {
	return r.DB.Create(user).Error
}