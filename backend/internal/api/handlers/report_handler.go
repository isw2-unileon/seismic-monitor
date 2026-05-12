package handlers

import (
	"net/http"
	"seismic-monitor/backend/internal/database"
	"seismic-monitor/backend/internal/models"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	Repo *database.ReportRepository
}

func (h *ReportHandler) HandleReport(c *gin.Context) {
	var report models.UserReport
	if err := c.ShouldBindJSON(&report); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Coordenadas inválidas"})
		return
	}

	count, err := h.Repo.RegisterReport(report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar reporte"})
		return
	}

	// Lógica de alerta temprana
	status := "Reporte recibido"
	if count >= 5 {
		status = "ALERTA: Múltiples reportes en tu zona. Posible sismo en curso."
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        status,
		"nearby_reports": count,
	})
}
