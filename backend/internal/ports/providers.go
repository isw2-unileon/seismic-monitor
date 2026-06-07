package ports

import "seismic-monitor/backend/internal/models"

// EarthquakeProvider is our PORT.
// Forces any adapter to return the structure our system understands.
type EarthquakeProvider interface {
	GetEarthquakes() (models.USGSResponse, error)
}
