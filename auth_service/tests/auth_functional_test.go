package tests

import (
	"testing"
)

// Skenario Functional / E2E Register -> Login
func TestE2E_RegisterAndLogin_RealDB(t *testing.T) {
	t.Log("Menjalankan Functional Test mengecek interaksi ke Database PostgreSQL...")

	// TODO: Hit API Register (http://localhost:8080/auth/register)
	// TODO: Cek apakah HTTP status 200 OK
	// TODO: Cek apakah data benar-benar masuk ke DB dengan SQL Query langsung

	// Langsung buat FAIL sesuai instruksi dosen karena DB belum siap
	t.Fatal("Functional Test Failed: Koneksi Database asli dan endpoint belum diimplementasi!")
}
