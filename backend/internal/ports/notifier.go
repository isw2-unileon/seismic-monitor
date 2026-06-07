package ports

import "seismic-monitor/backend/internal/models"

// NotificationService defines how the system sends alerts to the outside world.
type NotificationService interface {
	SendAlert(user models.User, sismo models.Feature) error
}
