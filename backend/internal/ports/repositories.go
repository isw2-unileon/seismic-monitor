package ports

import (
	"seismic-monitor/backend/internal/models"
	"time"
)

type EarthquakeRepository interface {
	GetEarthquakesSince(since time.Time) ([]models.Earthquake, error)
	GetFilteredEarthquakes(minMag float64, limit int) ([]models.Earthquake, error)
}
