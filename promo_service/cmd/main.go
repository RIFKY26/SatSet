package main

import (
	"fmt"
	"log"
	"net/http"

	"satset2/promo_service/domain"
	"satset2/promo_service/handler"
	"satset2/promo_service/repository"
	"satset2/promo_service/service"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=host.minikube.internal user=admin password=rahasia dbname=satset_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal connect ke database: %v", err)
	}

	// 1. Buat tabel promos
	err = db.AutoMigrate(&domain.Promo{})
	if err != nil {
		log.Fatalf("Gagal migrasi database: %v", err)
	}

	// 2. Seeding (Suntik promo otomatis untuk presentasi)
	var count int64
	db.Model(&domain.Promo{}).Where("promo_code = ?", "SATSET50").Count(&count)
	if count == 0 {
		db.Create(&domain.Promo{
			PromoID:       "P-001",
			PromoCode:     "SATSET50",
			MinOrderValue: 20000,
			MaxDiscount:   15000,
			DiscountPct:   50,
			Quota:         100,
		})
		fmt.Println("Promo 'SATSET50' berhasil disuntikkan ke database!")
	}

	repo := repository.NewSqlPromoRepository(db)
	promoService := service.NewPromoService(repo)
	promoHandler := handler.NewPromoHandler(promoService)

	http.HandleFunc("/promos/apply", promoHandler.ApplyPromoHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	fmt.Println("Promo Service berhasil terhubung ke PostgreSQL dan berjalan di port 8087...")
	log.Fatal(http.ListenAndServe(":8087", nil))
}
