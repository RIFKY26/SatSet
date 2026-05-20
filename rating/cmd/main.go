package main

import (
	"fmt"
	"satset2/rating/handler"
	"satset2/rating/repository"
	"satset2/rating/service"
)

func main() {
	fmt.Println("Menjalankan SatSet - Rating Service dengan Clean Architecture")

	repo := repository.NewRatingRepository()
	svc := service.NewRatingService(repo)
	h := handler.NewRatingHandler(svc)

	_ = h // Handler
}
