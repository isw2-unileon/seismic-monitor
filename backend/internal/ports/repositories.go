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

type UserRepository interface {
	CreateUser(user *models.User) error
	FindUserByEmail(email string) (*models.User, error)
	UpdateUserLocation(userID string, latitude, longitude, alertRadius, minMagnitude float64) error
	GetAffectedUsers(sismo models.Feature) ([]models.User, error)
	GetUsersNearLocation(lon, lat float64) ([]models.User, error)
}
