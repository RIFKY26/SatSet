package repository

import (
	"auth_service/domain"
	"gorm.io/gorm"
)

type AuthRepository interface {
	GetByEmail(email string) (*domain.User, error)
	CreateUser(user *domain.User) error
}

type SqlAuthRepository struct {
	DB *gorm.DB
}

func NewSqlAuthRepository(db *gorm.DB) AuthRepository {
	return &SqlAuthRepository{DB: db}
}

func (r *SqlAuthRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	// Mencari user berdasarkan email ke PostgreSQL
	result := r.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *SqlAuthRepository) CreateUser(user *domain.User) error {
	return r.DB.Create(user).Error
}