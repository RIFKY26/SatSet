package service

import (
	"testing"
	"user_service/domain"
	"user_service/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetByID_Unit_Success(t *testing.T) {
	// 1. Inisialisasi Gomock Controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 2. Buat objek MOCK dari repository (Tanpa Database Nyata)
	mockRepo := mocks.NewMockUserRepository(ctrl)

	// 3. Ajari Mock-nya: Kalau dipanggil "USR-001", kembalikan data dummy ini!
	expectedUser := &domain.User{
		ID:       "USR-001",
		Username: "budi_santoso",
		Email:    "budi.santoso@example.com",
		FullName: "Budi Santoso",
	}
	mockRepo.EXPECT().GetByID("USR-001").Return(expectedUser, nil)

	// 4. Masukkan mock repo ke dalam service
	svc := NewUserService(mockRepo)

	// 5. Jalankan fungsi yang dites
	result, err := svc.GetByID("USR-001")

	// 6. Validasi hasil (Pasti PASS karena diisolasi pakai Mock)
	assert.NoError(t, err)
	assert.Equal(t, "budi.santoso@example.com", result.Email)
	assert.Equal(t, "Budi Santoso", result.FullName)
}
