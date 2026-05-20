package service

import (
	"testing"
	"user_service/repository"
)

// Functional test mengakses real database
func TestGetByID_Functional_RealDB(t *testing.T) {
	t.Log("Menjalankan Functional Test - Mengakses Database MySQL Asli")

	// TODO: Setup koneksi database asli (gorm.DB / sql.DB)
	// db, err := sql.Open("mysql", "user:pass@tcp(localhost:3306)/satset_db")
	// realRepo := repository.NewUserRealRepository(db)

	// Karena koneksi DB dan real repository belum dibuat,
	// test ini dipastikan gagal sesuai instruksi / ekspektasi dosen.
	var realRepo repository.UserRepository = nil

	if realRepo == nil {
		t.Fatal("Functional Test Failed: Koneksi Database asli dan repository belum diimplementasi!")
	}

	// Jika suatu saat sudah ada DB, kodenya akan lanjut ke sini:
	svc := NewUserService(realRepo)
	_, err := svc.GetByID("2401248")

	if err != nil {
		t.Errorf("Error mengambil data dari DB: %v", err)
	}
}
