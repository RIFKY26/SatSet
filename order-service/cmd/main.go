package main

import (
	"fmt"
	"log"
	"net/http"

	"satset2/order-service/domain"
	"satset2/order-service/handler"
	"satset2/order-service/repository"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. DSN Koneksi ke Docker PostgreSQL
	dsn := "host=host.minikube.internal user=admin password=rahasia dbname=satset_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal connect ke database: %v", err)
	}

	// 2. AutoMigrate untuk membuat tabel orders otomatis
	err = db.AutoMigrate(&domain.Order{})
	if err != nil {
		log.Fatalf("Gagal migrasi database: %v", err)
	}

	// 3. Dependency Injection
	repo := repository.NewSqlOrderRepository(db)
	orderHandler := handler.NewOrderHandler(repo)

	// 4. Daftarkan Endpoint
	http.HandleFunc("/orders", orderHandler.CreateOrderHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	fmt.Println("Order Service berhasil terhubung ke PostgreSQL dan berjalan di port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
