package service

import (
	"auth_service/domain"
	"auth_service/mocks"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestLogin_Unit_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)

	mockRepo.EXPECT().
		GetByEmail("test@satset.com").
		Return(&domain.User{Email: "test@satset.com", Password: "hashedpassword"}, nil)

	svc := NewAuthService(mockRepo)

	_, err := svc.Login("test@satset.com", "password123")

	if err != nil {
		t.Errorf("Ekspektasi tidak ada error, tapi dapat: %v", err)
	}
}
