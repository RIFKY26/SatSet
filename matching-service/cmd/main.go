package main

import (
	"log"
	"net/http"

	"satset2/matching-service/handler"
)

func main() {
	http.HandleFunc("/match", handler.MatchDriverHandler)
	log.Println("Matching Service berjalan di :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
