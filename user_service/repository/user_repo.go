package repository

import "user_service/domain"

type UserRepository interface {
	GetByID(id string) (*domain.User, error)
	UpdateProfile(user *domain.User) error
	DeleteUser(id string) error
}
