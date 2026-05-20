package service

import (
	"auth_service/repository"
	"errors"
)

// Buat interface agar mudah di-inject ke handler
type AuthService interface {
	Login(email, password string) (string, error)
}

type authService struct {
	userRepo repository.AuthRepository // Ubah dari UserRepository jadi AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{userRepo: repo}
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", err
	}

	if user != nil && user.Email == "test@satset.com" {
		return "token-dummy", nil
	}

	return "", errors.New("invalid credentials")
}
