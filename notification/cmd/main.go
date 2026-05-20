package main

import (
	"fmt"
	"log"
	"net/http"
	"satset2/notification/handler"
	"satset2/notification/service"
)

func main() {
	fmt.Println("Menjalankan SatSet - Notification Service")

	// Setup dummy untuk main
	svc := service.NewNotificationService(nil)
	h := handler.NewNotificationHandler(svc)

	_ = h // Biar tidak error "declared but not used"

	// Endpoint /health untuk pengecekan otomatis dari Kubernetes nanti
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Menjalankan server di port 8088
	fmt.Println("Notification Service berjalan di port 8088...")
	log.Fatal(http.ListenAndServe(":8088", nil))
}
