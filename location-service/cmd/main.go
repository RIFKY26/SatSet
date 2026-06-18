package main

import (
	"fmt"
	"log"
	"net/http"

	"satset2/location-service/domain"
	"satset2/location-service/handler"
	"satset2/location-service/repository"
	"satset2/location-service/service"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=host.minikube.internal user=admin password=rahasia dbname=satset_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal connect ke database: %v", err)
	}

	err = db.AutoMigrate(&domain.DriverLocation{})
	if err != nil {
		log.Fatalf("Gagal migrasi database: %v", err)
	}

	// SEEDING: Memasukkan data supir palsu otomatis ke database
	var count int64
	db.Model(&domain.DriverLocation{}).Count(&count)
	if count == 0 {
		db.Create(&domain.DriverLocation{DriverID: "DRV-123", Latitude: -6.200100, Longitude: 106.816700}) // Sangat dekat dengan User
		db.Create(&domain.DriverLocation{DriverID: "DRV-999", Latitude: -6.500000, Longitude: 106.900000}) // Jauh
		fmt.Println("Berhasil menyuntikkan koordinat DRV-123 ke database!")
	}

	repo := repository.NewSqlLocationRepository(db)
	locSvc := service.NewLocationService(repo)
	locHandler := handler.NewLocationHandler(locSvc)

	http.HandleFunc("/location/update", locHandler.UpdateLocation)
	http.HandleFunc("/location/nearby", locHandler.GetNearbyDrivers)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	fmt.Println("Location Service terhubung ke DB dan berjalan di port 8083...")
	log.Fatal(http.ListenAndServe(":8083", nil))
}
