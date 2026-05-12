package handlers

import (
	"seismic-monitor/backend/internal/database"
	"seismic-monitor/backend/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	Repo       *database.ReportRepository
	UserRepo   *database.UserRepository
	AlertQueue chan<- models.AlertMessage // Canal para enviar a la cola
}

func (h *ReportHandler) HandleReport(c *gin.Context) {
	var report models.UserReport
	if err := c.ShouldBindJSON(&report); err != nil {
		c.JSON(400, gin.H{"error": "Datos inválidos"})
		return
	}

	// 1. Registrar y contar clúster
	count, err := h.Repo.RegisterReport(report)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error de base de datos"})
		return
	}

	// 2. Si hay 5 o más reportes, disparamos la alerta masiva
	if count == 5 { // Solo lo hacemos cuando llegamos exactamente a 5 para no spamear
		users, _ := h.UserRepo.GetUsersNearLocation(report.Longitude, report.Latitude)

		for _, user := range users {
			// Creamos un "Sismo ficticio" para la alerta comunitaria
			fakeSismo := models.Feature{
				ID:   "COMUNIDAD-" + time.Now().Format("20060102-150405"),
				Type: "Feature",
				Info: models.EarthquakeProps{
					Place: "Actividad reportada por usuarios en tu zona",
					Mag:   0.0, // Aún no hay magnitud oficial
					Time:  time.Now().UnixMilli(),
				},
				Geometry: models.EarthquakeGeometry{
					Type:        "Point",
					Coordinates: []float64{report.Longitude, report.Latitude, 0.0},
				},
			}

			// Enviamos a la cola que ya usa el NotificationWorker
			h.AlertQueue <- models.AlertMessage{
				User:  user,
				Sismo: fakeSismo,
			}
		}
	}

	c.JSON(200, gin.H{"status": "ok", "nearby": count})
}
