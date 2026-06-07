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
	// Use gen_random_uuid() for the ID and ST_MakePoint for the location
	insertQuery := `
		INSERT INTO reported_earthquakes (reported_earthquake_id, location, reported_at)
		VALUES (gen_random_uuid(), ST_SetSRID(ST_MakePoint($1, $2), 4326), NOW())
		RETURNING reported_at`

	err := r.DB.QueryRow(insertQuery, report.Longitude, report.Latitude).Scan(&report.ReportedAt)
	if err != nil {
		return 0, fmt.Errorf("error al guardar reporte: %w", err)
	}

	// Calculate the cluster: How many people have reported within a 30km radius
	// in the last 10 minutes?
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


func (r *ReportRepository) CleanOldReports(interval string) (int64, error) {
	query := fmt.Sprintf("DELETE FROM reported_earthquakes WHERE reported_at < NOW() - INTERVAL '%s'", interval)

	result, err := r.DB.Exec(query)
	if err != nil {
		return 0, fmt.Errorf("error limpiando reportes antiguos: %w", err)
	}

	return result.RowsAffected()
}
