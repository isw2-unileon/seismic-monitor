package ports

import "seismic-monitor/backend/internal/models"

// EarthquakeProvider es nuestro PUERTO.
// Obliga a cualquier adaptador a devolver la estructura que nuestro sistema entiende.
type EarthquakeProvider interface {
	GetEarthquakes() (models.USGSResponse, error)
}
