package main

import (
	"fmt"
	"log"
	"net/http"
	"satset2/payment-service/repository"
	"satset2/payment-service/service"
)

func main() {
	fmt.Println("Payment Service is running...")

	repo := repository.NewPaymentRepository()
	paymentSvc := service.NewPaymentService(repo)

	// Logika bawaanmu tetap kita simpan agar tidak rusak
	tx, err := paymentSvc.ProcessPayment("ORD-001", "USR-01", 150000, "authorize", "unique-1")
	if err == nil {
		fmt.Printf("Transaction Created: %s dengan status %s\n", tx.TransactionID, tx.Status)
	}

	// Endpoint /health untuk Kubernetes
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Menjalankan server di port 8089
	fmt.Println("Payment Service berjalan di port 8089...")
	log.Fatal(http.ListenAndServe(":8089", nil))
}
