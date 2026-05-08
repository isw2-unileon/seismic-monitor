package ports

import (
	"seismic-monitor/backend/internal/models"
	"time"
)

type EarthquakeRepository interface {
	GetEarthquakesSince(since time.Time) ([]models.Feature, error)
	GetFilteredEarthquakes(minMag float64, limit int) ([]models.Feature, error)
	SaveEarthquake(eq models.Feature) error
}
