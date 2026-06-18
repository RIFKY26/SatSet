package service

import (
	"auth_service/repository"
	"errors"
)

type AuthService interface {
	Login(email, password string) (string, error)
}

type authService struct {
	userRepo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{userRepo: repo}
}

func (s *authService) Login(email, password string) (string, error) {
	// 1. Cari user di Database asli
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", errors.New("email tidak terdaftar")
	}

	// 2. Cek password (Di production pakai bcrypt, untuk tugas ini string biasa dulu)
	if user.Password == password {
		return "token-jwt-asli-nanti-disini", nil
	}

	return "", errors.New("password salah")
}