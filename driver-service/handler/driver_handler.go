package handler

import (
	"encoding/json"
	"net/http"

	"satset2/driver-service/service"
)

type DriverHandler struct {
	service *service.DriverService
}

func NewDriverHandler(svc *service.DriverService) *DriverHandler {
	return &DriverHandler{service: svc}
}

type registerRequest struct {
	DriverID string `json:"driver_id"`
}
type assignRequest struct {
	DriverID string `json:"driver_id"`
}
type completeRequest struct {
	DriverID string `json:"driver_id"`
}
type statusRequest struct {
	DriverID string `json:"driver_id"`
}

type statusResponse struct {
	DriverID           string `json:"driver_id"`
	ConnectionStatus   string `json:"connection_status"`
	AvailabilityStatus string `json:"availability_status"`
}

func (h *DriverHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.DriverID == "" {
		respondError(w, http.StatusBadRequest, "driver_id is required")
		return
	}
	if err := h.service.RegisterDriver(req.DriverID); err != nil {
		respondError(w, http.StatusConflict, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, successResponse{Message: "driver registered"})
}

func (h *DriverHandler) Assign(w http.ResponseWriter, r *http.Request) {
	var req assignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.DriverID == "" {
		respondError(w, http.StatusBadRequest, "driver_id is required")
		return
	}
	if err := h.service.AssignOrder(req.DriverID); err != nil {
		respondError(w, conflictOrNotFound(err), err.Error())
		return
	}
	respondJSON(w, http.StatusOK, successResponse{Message: "order assigned"})
}

func (h *DriverHandler) Complete(w http.ResponseWriter, r *http.Request) {
	var req completeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.DriverID == "" {
		respondError(w, http.StatusBadRequest, "driver_id is required")
		return
	}
	if err := h.service.CompleteOrder(req.DriverID); err != nil {
		respondError(w, conflictOrNotFound(err), err.Error())
		return
	}
	respondJSON(w, http.StatusOK, successResponse{Message: "order completed"})
}

func (h *DriverHandler) Status(w http.ResponseWriter, r *http.Request) {
	var req statusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.DriverID == "" {
		respondError(w, http.StatusBadRequest, "driver_id is required")
		return
	}
	d, err := h.service.GetDriver(req.DriverID)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, statusResponse{
		DriverID:           d.ID,
		ConnectionStatus:   string(d.ConnectionStatus),
		AvailabilityStatus: string(d.AvailabilityStatus),
	})
}
