package main

import (
	"auth_service/handler"
	"auth_service/repository"
	"auth_service/service"
	"fmt"
	"log"
	"net/http"
)

func main() {
	var repo repository.AuthRepository
	authService := service.NewAuthService(repo)
	authHandler := handler.NewAuthHandler(authService)

	http.HandleFunc("/login", authHandler.Login)

	// Tambahan endpoint health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	fmt.Println("Auth Service berjalan di port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
