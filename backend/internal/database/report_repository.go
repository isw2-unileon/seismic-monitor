package database

import (
	"database/sql"
	"fmt"
	"seismic-monitor/backend/internal/models"
)

type ReportRepository struct {
	DB *sql.DB
}

func (r *ReportRepository) RegisterReport(report models.UserReport) (int, error) {
	// 1. Insertamos el reporte usando la estructura de tu amigo
	// Usamos gen_random_uuid() para el ID y ST_MakePoint para la ubicación
	insertQuery := `
		INSERT INTO reported_earthquakes (reported_earthquake_id, location, reported_at)
		VALUES (gen_random_uuid(), ST_SetSRID(ST_MakePoint($1, $2), 4326), NOW())
		RETURNING reported_at`

	err := r.DB.QueryRow(insertQuery, report.Longitude, report.Latitude).Scan(&report.ReportedAt)
	if err != nil {
		return 0, fmt.Errorf("error al guardar reporte: %w", err)
	}

	// 2. Calculamos el clúster: ¿Cuánta gente ha reportado en un radio de 30km
	// en los últimos 10 minutos?
	countQuery := `
		SELECT COUNT(*) 
		FROM reported_earthquakes 
		WHERE ST_DWithin(location, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, 30000)
		  AND reported_at > NOW() - INTERVAL '2 minutes'`

	var count int
	err = r.DB.QueryRow(countQuery, report.Longitude, report.Latitude).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error al contar clúster: %w", err)
	}

	return count, nil
}
