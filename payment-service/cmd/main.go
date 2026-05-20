package main

import (
	"fmt"
	"satset2/payment-service/repository"
	"satset2/payment-service/service"
)

func main() {
	fmt.Println("Payment Service is running...")

	repo := repository.NewPaymentRepository()
	paymentSvc := service.NewPaymentService(repo)

	tx, err := paymentSvc.ProcessPayment("ORD-001", "USR-01", 150000, "authorize", "unique-1")
	if err == nil {
		fmt.Printf("Transaction Created: %s dengan status %s\n", tx.TransactionID, tx.Status)
	}
}
