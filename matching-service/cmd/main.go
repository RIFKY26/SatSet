package main

import (
	"fmt"
	"log"
	"net/http"

	"satset2/matching-service/client"
	"satset2/matching-service/handler"
	"satset2/matching-service/service"
)

func main() {
	// 1. Siapkan kabel untuk menelepon location-service di port 8083
	locClient := client.NewLocationHTTPClient("http://location-service:80")

	// 2. Dependency Injection
	matchSvc := service.NewMatchService(locClient)
	matchHandler := handler.NewMatchHandler(matchSvc)

	// 3. Daftarkan Endpoint
	http.HandleFunc("/match", matchHandler.MatchDriverHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	fmt.Println("Matching Service berjalan di port 8084...")
	log.Fatal(http.ListenAndServe(":8084", nil))
}
