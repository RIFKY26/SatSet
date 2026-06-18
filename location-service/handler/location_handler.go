package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"satset2/location-service/domain"
	"satset2/location-service/service"
)

type LocationHandler struct {
	locService *service.LocationService
}

func NewLocationHandler(s *service.LocationService) *LocationHandler {
	return &LocationHandler{locService: s}
}

func (h *LocationHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	var req domain.DriverLocation
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	h.locService.UpdateDriverLocation(&req)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "updated"}`))
}

func (h *LocationHandler) GetNearbyDrivers(w http.ResponseWriter, r *http.Request) {
	lat, _ := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	lng, _ := strconv.ParseFloat(r.URL.Query().Get("lng"), 64)
	radius, _ := strconv.ParseFloat(r.URL.Query().Get("radius"), 64)

	drivers, err := h.locService.GetNearestDrivers(lat, lng, radius)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drivers)
}
