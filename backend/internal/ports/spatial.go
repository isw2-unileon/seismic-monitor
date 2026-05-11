package ports

import "seismic-monitor/backend/internal/models"

// SpatialRepository define la interfaz que debe cumplir cualquier
// sistema que busque usuarios afectados por un sismo.
type SpatialRepository interface {
	GetAffectedUsers(sismo models.Feature) ([]models.User, error)
}
