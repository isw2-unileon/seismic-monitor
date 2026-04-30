package usgs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"seismic-monitor/backend/internal/models"
)

// USGSAdapter es el adaptador que se conecta al exterior.
type USGSAdapter struct {
	URL string
}

// GetEarthquakes cumple con la interfaz ports.EarthquakeProvider
func (a *USGSAdapter) GetEarthquakes() (models.USGSResponse, error) {
	data, err := FetchData(a.URL)
	if err != nil {
		return models.USGSResponse{}, fmt.Errorf("fallo en fetch: %w", err)
	}

	response, err := ParseUSGSData(data)
	if err != nil {
		return models.USGSResponse{}, fmt.Errorf("fallo en parseo: %w", err)
	}

	return response, nil
}

// FetchData realiza la petición HTTP GET a la URL proporcionada
func FetchData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// ParseUSGSData deserializa el JSON en la estructura USGSResponse
func ParseUSGSData(data []byte) (models.USGSResponse, error) {
	var response models.USGSResponse
	err := json.Unmarshal(data, &response)
	return response, err
}
