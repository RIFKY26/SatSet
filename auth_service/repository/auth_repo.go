package repository

import "auth_service/domain"

type AuthRepository interface {
	GetByEmail(email string) (*domain.User, error)
	CreateUser(user *domain.User) error
}
