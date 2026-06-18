package main

import (
	"fmt"
	"log"
	"net/http"

	"satset2/driver-service/domain"
	"satset2/driver-service/handler"
	"satset2/driver-service/repository"
	"satset2/driver-service/service"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. DSN Koneksi ke Docker PostgreSQL
	dsn := "host=localhost user=admin password=rahasia dbname=satset_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal connect ke database: %v", err)
	}

	// 2. AutoMigrate untuk membuat tabel drivers otomatis
	err = db.AutoMigrate(&domain.Driver{})
	if err != nil {
		log.Fatalf("Gagal migrasi database: %v", err)
	}

	// 3. Dependency Injection
	repo := repository.NewSqlDriverRepository(db)
	svc := service.NewDriverService(repo)
	h := handler.NewDriverHandler(svc)

	// 4. Daftarkan Endpoint
	mux := http.NewServeMux()
	mux.HandleFunc("/drivers/register", h.Register)
	mux.HandleFunc("/drivers/assign", h.Assign)
	mux.HandleFunc("/drivers/complete", h.Complete)
	mux.HandleFunc("/drivers/status", h.Status)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	fmt.Println("Driver Service berhasil terhubung ke PostgreSQL dan berjalan di port 8082...")
	log.Fatal(http.ListenAndServe(":8082", mux))
}