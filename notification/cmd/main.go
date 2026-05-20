package main

import (
	"fmt"
	"satset2/notification/handler"
	"satset2/notification/service"
)

func main() {
	fmt.Println("Menjalankan SatSet - Notification Service")

	// Setup dummy untuk main
	svc := service.NewNotificationService(nil)
	h := handler.NewNotificationHandler(svc)

	_ = h // Biar tidak error "declared but not used"
}
