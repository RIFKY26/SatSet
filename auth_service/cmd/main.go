package main

import (
	"auth_service/handler"
	"auth_service/repository"
	"auth_service/service"
	"fmt"
	"net/http"
)

func main() {
	var repo repository.AuthRepository
	authService := service.NewAuthService(repo)
	authHandler := handler.NewAuthHandler(authService)

	http.HandleFunc("/login", authHandler.Login)
	fmt.Println("Auth Service berjalan di port 8081...")
	http.ListenAndServe(":8081", nil)
}
