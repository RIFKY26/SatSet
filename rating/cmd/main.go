package main

import (
	"fmt"
	"log"
	"net/http"

	"satset2/rating/domain"
	"satset2/rating/handler"
	"satset2/rating/repository"
	"satset2/rating/service"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=host.minikube.internal user=admin password=rahasia dbname=satset_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal connect ke database: %v", err)
	}

	// AutoMigrate untuk membuat tabel ratings otomatis
	err = db.AutoMigrate(&domain.Rating{})
	if err != nil {
		log.Fatalf("Gagal migrasi database: %v", err)
	}

	repo := repository.NewSqlRatingRepository(db)
	ratingSvc := service.NewRatingService(repo)
	ratingHandler := handler.NewRatingHandler(ratingSvc)

	http.HandleFunc("/ratings", ratingHandler.SubmitRatingHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	fmt.Println("Rating Service berhasil terhubung ke PostgreSQL dan berjalan di port 8088...")
	log.Fatal(http.ListenAndServe(":8088", nil))
}
