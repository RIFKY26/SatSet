package main

import (
	"log"
	"net/http"
	"satset2/promo-service/handler"
	"satset2/promo-service/repository"
	"satset2/promo-service/service"
	"time"
)

func main() {
	// Perakitan Komponen (Dependency Injection)
	repo := repository.NewPromoRepository()
	promoSvc := service.NewPromoService(repo)
	promoHandler := handler.NewPromoHandler(promoSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/promo/apply", promoHandler.HandleApplyPromo)

	// TAMBAHAN WAJIB: Endpoint health agar Kubernetes bisa mengecek status Promo Service
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:         ":8086",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("Promo Service starting on :8086...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
