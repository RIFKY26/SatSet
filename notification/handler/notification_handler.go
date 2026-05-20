package handler

import (
	"net/http"
	"satset2/notification/service"
)

type NotificationHandler struct {
	NotifyService *service.NotificationService
}

func NewNotificationHandler(svc *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{NotifyService: svc}
}

func (h *NotificationHandler) HandleSendNotification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"message": "TDD: handler HTTP belum diimplementasikan"}`))
}
