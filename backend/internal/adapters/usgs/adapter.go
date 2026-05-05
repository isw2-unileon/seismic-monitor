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
	data, err := a.fetchData(a.URL)
	if err != nil {
		return models.USGSResponse{}, fmt.Errorf("fallo en fetch: %w", err)
	}

	response, err := a.parseUSGSData(data)
	if err != nil {
		return models.USGSResponse{}, fmt.Errorf("fallo en parseo: %w", err)
	}

	return response, nil
}

func (a *USGSAdapter) fetchData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (a *USGSAdapter) parseUSGSData(data []byte) (models.USGSResponse, error) {
	var response models.USGSResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return models.USGSResponse{}, err
	}
	return response, nil
}
