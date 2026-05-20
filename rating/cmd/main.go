package main

import (
	"fmt"
	"log"
	"net/http"
	"satset2/rating/handler"
	"satset2/rating/repository"
	"satset2/rating/service"
)

func main() {
	fmt.Println("Menjalankan SatSet - Rating Service dengan Clean Architecture")

	repo := repository.NewRatingRepository()
	svc := service.NewRatingService(repo)
	h := handler.NewRatingHandler(svc)

	_ = h // Biar variabel h tidak error 'declared but not used' jika belum ada route-nya

	// Endpoint /health untuk pengecekan otomatis dari Kubernetes nanti
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Menjalankan server di port 8087
	fmt.Println("Rating Service berjalan di port 8087...")
	log.Fatal(http.ListenAndServe(":8087", nil))
}
