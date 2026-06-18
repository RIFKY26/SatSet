package main

import (
	"fmt"
	"log"
	"net/http"

	"auth_service/domain"
	"auth_service/handler"
	"auth_service/repository"
	"auth_service/service"

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

	// 2. AutoMigrate untuk membuat tabel auth_users
	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		log.Fatalf("Gagal migrasi database: %v", err)
	}

	// 3. Dependency Injection
	repo := repository.NewSqlAuthRepository(db)
	authSvc := service.NewAuthService(repo)
	authHandler := handler.NewAuthHandler(authSvc)

	// 4. Setup Router standar bawaan Go
	http.HandleFunc("/login", authHandler.Login)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	fmt.Println("Auth Service berhasil terhubung ke PostgreSQL dan berjalan di port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
