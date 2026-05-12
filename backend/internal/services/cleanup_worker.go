package services

import (
	"log/slog"
	"seismic-monitor/backend/internal/database"
	"time"
)

// StartReportCleanupWorker ejecuta una limpieza cada 30 minutos
func StartReportCleanupWorker(repo *database.ReportRepository) {
	ticker := time.NewTicker(30 * time.Minute)

	go func() {
		for range ticker.C {
			slog.Info("[Cleanup Worker] Iniciando limpieza de reportes antiguos...")

			deleted, err := repo.CleanOldReports("1 hour")
			if err != nil {
				slog.Error("[Cleanup Worker] Fallo en la limpieza", "error", err)
				continue
			}

			if deleted > 0 {
				slog.Info("[Cleanup Worker] Limpieza completada", "filas_eliminadas", deleted)
			}
		}
	}()
}
