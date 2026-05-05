package ports

import "seismic-monitor/backend/internal/models"

// SpatialRepository define la función que el Desarrollador 1 implementará
type SpatialRepository interface {
	GetAffectedUsers(sismo models.Feature) ([]models.User, error)
}
