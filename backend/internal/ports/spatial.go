package ports

import "seismic-monitor/backend/internal/models"

type SpatialRepository interface {
	GetAffectedUsers(sismo models.Feature) ([]models.User, error)
}
