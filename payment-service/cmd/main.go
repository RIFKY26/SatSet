package main

import (
	"fmt"
	"log"
	"net/http"

	"satset2/payment-service/domain"
	"satset2/payment-service/handler"
	"satset2/payment-service/repository"
	"satset2/payment-service/service"

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

	// 2. AutoMigrate untuk membuat tabel transactions otomatis
	err = db.AutoMigrate(&domain.Transaction{})
	if err != nil {
		log.Fatalf("Gagal migrasi database: %v", err)
	}

	// 3. Dependency Injection
	repo := repository.NewSqlPaymentRepository(db)
	paymentSvc := service.NewPaymentService(repo)
	paymentHandler := handler.NewPaymentHandler(paymentSvc)

	// 4. Daftarkan Endpoint HTTP
	http.HandleFunc("/payments", paymentHandler.ProcessPaymentHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Server berjalan di port 8086
	fmt.Println("Payment Service berhasil terhubung ke PostgreSQL dan berjalan di port 8086...")
	log.Fatal(http.ListenAndServe(":8086", nil))
}
