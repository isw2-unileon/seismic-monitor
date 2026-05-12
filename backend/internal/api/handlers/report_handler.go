package handlers

import (
	"net/http"
	"sync"
	"time"

	"seismic-monitor/backend/internal/database"
	"seismic-monitor/backend/internal/models"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	Repo       *database.ReportRepository
	UserRepo   *database.UserRepository
	AlertQueue chan<- models.AlertMessage

	// Sistema anti-spam
	lastReports sync.Map // Guarda IP -> time.Time
	limit       time.Duration
}

func NewReportHandler(repo *database.ReportRepository, userRepo *database.UserRepository, queue chan<- models.AlertMessage) *ReportHandler {
	return &ReportHandler{
		Repo:       repo,
		UserRepo:   userRepo,
		AlertQueue: queue,
		limit:      2 * time.Minute, // Un reporte cada 2 minutos por IP
	}
}

func (h *ReportHandler) HandleReport(c *gin.Context) {
	userIP := c.ClientIP()

	// 1. Verificar si la IP está en "enfriamiento"
	if lastTime, ok := h.lastReports.Load(userIP); ok {
		if time.Since(lastTime.(time.Time)) < h.limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Has enviado un reporte recientemente. Por favor, espera un poco.",
			})
			return
		}
	}

	var report models.UserReport
	if err := c.ShouldBindJSON(&report); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// 2. Registrar el reporte y obtener el conteo de clúster
	count, err := h.Repo.RegisterReport(report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error de base de datos"})
		return
	}

	// 3. Si todo ok, actualizamos el timestamp del usuario
	h.lastReports.Store(userIP, time.Now())

	// 4. Lógica de alerta masiva (solo si llegamos al umbral de 5)
	if count == 5 {
		// ... (mismo código de búsqueda de usuarios y envío a la AlertQueue que ya tenemos)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reporte recibido correctamente",
		"nearby":  count,
	})
}
