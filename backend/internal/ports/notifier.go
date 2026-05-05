package ports

import "seismic-monitor/backend/internal/models"

// NotificationService define cómo el sistema envía alertas al mundo exterior.
type NotificationService interface {
	SendAlert(user models.User, sismo models.Feature) error
}
