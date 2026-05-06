package main

import (
	"log"
	"net/http"

	"satset2/location-service/client"
	"satset2/location-service/handler"
	"satset2/location-service/service"
)

func main() {
	// 1. Inisialisasi HTTP Client untuk menembak driver-service
	driverClient := client.NewDriverHTTPClient("http://localhost:8082")

	// 2. Masukkan client ke dalam LocationService
	svc := service.NewLocationService(driverClient)

	// 3. Masukkan service ke dalam Handler
	h := handler.NewLocationHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/location/update", h.Update)
	mux.HandleFunc("/location/nearby", h.Nearby)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Location Service berjalan di :8083")
	log.Fatal(http.ListenAndServe(":8083", mux))
}
