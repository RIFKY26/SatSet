package tests

import (
	"testing"
)

func TestSubmitRating_FunctionalTest(t *testing.T) {
	t.Log("Menjalankan Functional Test mengecek interaksi ke Database...")

	// Kosong tanpa mock DB
	// err := svc.SubmitRating("ORD-REAL-1", "DRV-REAL-1", 5)
	// assert.NoError(t, err, "Harusnya bisa menyimpan rating langsung ke Database asli")

	// Langsung buat FAIL sesuai instruksi dosen karena DB belum siap
	t.Fatal("Functional Test Failed: Koneksi Database asli dan repository belum diimplementasi!")
}
