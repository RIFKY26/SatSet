package main

import (
	"fmt"
	"log"
	"net/http"

	"user_service/domain"
	"user_service/handler"
	"user_service/repository"
	"user_service/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Data koneksi (DSN) sesuai dengan yang ada di Docker
	dsn := "host=host.minikube.internal user=admin password=rahasia dbname=satset_db port=5432 sslmode=disable"

	// 2. Buka koneksi ke PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal koneksi ke database: %v", err)
	}

	// 3. AutoMigrate: Otomatis membuat tabel users
	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		log.Fatalf("Gagal migrasi database: %v", err)
	}

	// 4. Rakit aplikasinya (Dependency Injection)
	repo := repository.NewSqlUserRepository(db)
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)

	// 5. Setup Router (karena handler kamu menggunakan Gin)
	r := gin.Default()

	// Endpoint untuk ambil profile user
	r.GET("/users/:id", h.GetProfile)

	// Endpoint Health untuk mengecek Kubernetes nanti
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	fmt.Println("User Service berhasil terhubung ke PostgreSQL dan berjalan di port 8085...")
	// Menjalankan server di port 8085
	err = r.Run(":8085")
	if err != nil {
		log.Fatalf("Server gagal berjalan: %v", err)
	}
}
