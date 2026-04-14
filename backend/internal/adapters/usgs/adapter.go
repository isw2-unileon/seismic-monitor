package usgs

import (
	"fmt"

	"seismic-monitor/internal/models"
)

// USGSAdapter es el adaptador que se conecta al exterior.
type USGSAdapter struct {
	URL string
}

// GetEarthquakes cumple con la interfaz ports.EarthquakeProvider
func (a *USGSAdapter) GetEarthquakes() (models.USGSResponse, error) {
	// 1. Usamos tu función FetchData intacta
	data, err := FetchData(a.URL)
	if err != nil {
		return models.USGSResponse{}, fmt.Errorf("fallo en fetch: %w", err)
	}

	// 2. Usamos tu función ParseUSGSData intacta
	response, err := ParseUSGSData(data)
	if err != nil {
		return models.USGSResponse{}, fmt.Errorf("fallo en parseo: %w", err)
	}

	return response, nil
}
