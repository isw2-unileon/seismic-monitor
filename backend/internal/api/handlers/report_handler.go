package handlers

import (
	"net/http"
	"sync"
	"time"

	"seismic-monitor/backend/internal/database"
	"seismic-monitor/backend/internal/models"

	"github.com/gin-gonic/gin"
)

type ReportRepository interface {
	RegisterReport(report models.UserReport) (int, error)
}

type ReportHandler struct {
	Repo        ReportRepository         // 2. CHANGED: Previously *database.ReportRepository, now uses the interface
	UserRepo    *database.UserRepository // (If desired, in the future this can also use an interface)
	AlertQueue  chan<- models.AlertMessage
	lastReports sync.Map
	limit       time.Duration
}

// 3. CHANGED: The first argument now receives the 'ReportRepository' interface
func NewReportHandler(repo ReportRepository, userRepo *database.UserRepository, queue chan<- models.AlertMessage) *ReportHandler {
	return &ReportHandler{
		Repo:       repo,
		UserRepo:   userRepo,
		AlertQueue: queue,
		limit:      2 * time.Minute,
	}
}

func (h *ReportHandler) HandleReport(c *gin.Context) {
	userIP := c.ClientIP()

	// Check if the IP is in "cooldown"
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

	if report.Latitude < -90 || report.Latitude > 90 || report.Longitude < -180 || report.Longitude > 180 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Coordenadas geográficas fuera de los límites de la Tierra",
		})
		return
	}

	count, err := h.Repo.RegisterReport(report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error de base de datos"})
		return
	}

	h.lastReports.Store(userIP, time.Now())

	// Mass alert logic (only if we reach the threshold of 5)
	if count == 5 {
		// ... (same user search and AlertQueue sending code as before)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reporte recibido correctamente",
		"nearby":  count,
	})
}
