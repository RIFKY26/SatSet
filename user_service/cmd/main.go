package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("User Service Satset siap melayani!")

	// Kita tambahkan endpoint /health (sangat berguna untuk dicek oleh Kubernetes nanti)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Menjalankan server di port 8085
	fmt.Println("User Service berjalan di port 8085...")
	log.Fatal(http.ListenAndServe(":8085", nil))
}
