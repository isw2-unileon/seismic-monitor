package services

import (
	"seismic-monitor/backend/internal/models"
	"seismic-monitor/backend/internal/ports"
	"time"
)

type EarthquakeService struct {
	repo ports.EarthquakeRepository
}

func NewEarthquakeService(repo ports.EarthquakeRepository) *EarthquakeService {
	return &EarthquakeService{repo: repo}
}

func (s *EarthquakeService) GetRecentEarthquakes() ([]models.Earthquake, error) {
	// Business logic: calculate exactly 1 hour ago
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	return s.repo.GetEarthquakesSince(oneHourAgo)
}

func (s *EarthquakeService) GetHistory() ([]models.Earthquake, error) {
	// Return earthquakes from the last hour for history
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	return s.repo.GetEarthquakesSince(oneHourAgo)
}

func (s *EarthquakeService) GetFilteredEarthquakes(minMag float64, limit int) ([]models.Earthquake, error) {
	return s.repo.GetFilteredEarthquakes(minMag, limit)
}
