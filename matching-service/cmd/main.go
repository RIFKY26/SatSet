package main

import (
	"log"
	"net/http"
	"satset2/matching-service/handler"
)

func main() {
	http.HandleFunc("/match", handler.MatchDriverHandler)

	// Tambahan endpoint health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Matching Service berjalan di port 8084...")
	log.Fatal(http.ListenAndServe(":8084", nil))
}
