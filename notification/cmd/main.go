package main

import (
	"fmt"
	"log"
	"net/http"

	"satset2/notification/domain"
	"satset2/notification/handler"
	"satset2/notification/repository"
	"satset2/notification/service"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=admin password=rahasia dbname=satset_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal connect ke database: %v", err)
	}

	err = db.AutoMigrate(&domain.Notification{})
	if err != nil {
		log.Fatalf("Gagal migrasi database: %v", err)
	}

	repo := repository.NewSqlNotificationRepository(db)
	notifSvc := service.NewNotificationService(repo)
	notifHandler := handler.NewNotificationHandler(notifSvc)

	http.HandleFunc("/notifications/send", notifHandler.SendNotificationHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	fmt.Println("Notification Service berhasil terhubung ke PostgreSQL dan berjalan di port 8089...")
	log.Fatal(http.ListenAndServe(":8089", nil))
}
