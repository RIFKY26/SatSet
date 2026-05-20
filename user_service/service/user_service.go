package service

import (
	"user_service/domain"
	"user_service/repository"
)

type UserService interface {
	GetByID(id string) (*domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) GetByID(id string) (*domain.User, error) {
	return s.repo.GetByID(id)
}
