package handler

import (
	"encoding/json"
	"net/http"
	"satset2/location-service/service"
	"strconv"
	"time"
)

type LocationHandler struct {
	service *service.LocationService
}

func NewLocationHandler(svc *service.LocationService) *LocationHandler {
	return &LocationHandler{service: svc}
}

type updateLocationRequest struct {
	DriverID  string  `json:"driver_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (h *LocationHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req updateLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.DriverID == "" {
		respondError(w, http.StatusBadRequest, "driver_id is required")
		return
	}

	h.service.UpdateLocation(req.DriverID, req.Latitude, req.Longitude, time.Now().Unix())
	respondJSON(w, http.StatusOK, successResponse{Message: "location updated"})
}

func (h *LocationHandler) Nearby(w http.ResponseWriter, r *http.Request) {
	lat, err := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "lat must be a float")
		return
	}
	lng, err := strconv.ParseFloat(r.URL.Query().Get("lng"), 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "lng must be a float")
		return
	}
	radius, err := strconv.ParseFloat(r.URL.Query().Get("radius"), 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "radius must be a float")
		return
	}

	drivers := h.service.GetNearbyDrivers(lat, lng, radius)
	respondJSON(w, http.StatusOK, drivers)
}
