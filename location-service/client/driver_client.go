package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DriverHTTPClient struct {
	baseURL string
}

func NewDriverHTTPClient(baseURL string) *DriverHTTPClient {
	return &DriverHTTPClient{baseURL: baseURL}
}

func (c *DriverHTTPClient) IsDriverAvailable(driverID string) (bool, error) {
	url := fmt.Sprintf("%s/drivers/status", c.baseURL)
	reqBody := map[string]string{"driver_id": driverID}
	bodyBytes, _ := json.Marshal(reqBody)

	resp, err := http.Post(url, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return false, fmt.Errorf("gagal menghubungi driver-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var statusResp struct {
		ConnectionStatus   string `json:"connection_status"`
		AvailabilityStatus string `json:"availability_status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		return false, fmt.Errorf("gagal decode respons: %w", err)
	}

	return statusResp.ConnectionStatus == "ONLINE" && statusResp.AvailabilityStatus == "AVAILABLE", nil
}
