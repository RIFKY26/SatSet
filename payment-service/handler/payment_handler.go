package handler

import (
	"net/http"
	"satset2/payment-service/domain"
)

type PaymentHandler struct {
	service domain.PaymentService
}

func NewPaymentHandler(s domain.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: s}
}

func (h *PaymentHandler) HandleProcessPayment(w http.ResponseWriter, r *http.Request) {
	// Logika parsing HTTP request dan pemanggilan service
}
