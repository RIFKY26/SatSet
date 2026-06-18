package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LocationHTTPClient struct {
	baseURL string
}

func NewLocationHTTPClient(url string) *LocationHTTPClient {
	return &LocationHTTPClient{baseURL: url}
}

func (c *LocationHTTPClient) GetNearestDriver(lat, lng float64) (string, error) {
	// Nembak API Location Service yang berjalan di port 8083 dengan radius 10 km
	url := fmt.Sprintf("%s/location/nearby?lat=%f&lng=%f&radius=10", c.baseURL, lat, lng)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("gagal memanggil location-service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gagal mendapatkan lokasi")
	}

	// Menangkap balasan dari location-service
	var drivers []struct {
		DriverID string  `json:"driver_id"`
		Distance float64 `json:"distance"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&drivers); err != nil {
		return "", err
	}

	if len(drivers) == 0 {
		return "", fmt.Errorf("tidak ada driver yang tersedia di sekitarmu")
	}

	// Kembalikan Driver ID urutan pertama (yang paling dekat)
	return drivers[0].DriverID, nil
}