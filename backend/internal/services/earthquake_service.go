 package services

      import (
         "seismic-monitor/backend/internal/models"
         "seismic-monitor/backend/internal/ports"
         "time"
     )

     type EarthquakeService struct {
         repo ports.EarthquakeRepository
     }

     func NewEarthquakeService(repo ports.EarthquakeRepository)
      	*EarthquakeService {
         return &EarthquakeService{repo: repo}
     }

     func (s *EarthquakeService) GetRecentEarthquakes()
      ([]models.Earthquake, error) {
         // Business logic: calculate exactly 1 hour ago
         oneHourAgo := time.Now().Add(-1 * time.Hour)
         return s.repo.GetEarthquakesSince(oneHourAgo)
     }