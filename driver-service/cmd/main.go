package main

import (
	"log"
	"net/http"

	"satset2/driver-service/handler"
	"satset2/driver-service/service"
)

func main() {
	svc := service.NewDriverService()
	h := handler.NewDriverHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/drivers/register", h.Register)
	mux.HandleFunc("/drivers/assign", h.Assign)
	mux.HandleFunc("/drivers/complete", h.Complete)
	mux.HandleFunc("/drivers/status", h.Status)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Driver Service berjalan di :8082")
	log.Fatal(http.ListenAndServe(":8082", mux))
}
