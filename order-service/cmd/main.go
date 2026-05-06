package main

import (
	"log"
	"net/http"

	"satset2/order-service/handler" // Import handler yang baru dibuat
)

func main() {
	// Register handler dengan memanggil package handler
	http.HandleFunc("/orders", handler.CreateOrderHandler)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Order Service berjalan di :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
