package database

import (
	"database/sql"
	"fmt"
	"seismic-monitor/backend/internal/models"
)

type EarthquakeRepository struct {
	DB *sql.DB
}

func NewEarthquakeRepository(db *sql.DB) *EarthquakeRepository {
	return &EarthquakeRepository{DB: db}
}

// GetFilteredEarthquakes obtiene sismos filtrados por magnitud mínima y límite
func (r *EarthquakeRepository) GetFilteredEarthquakes(minMag float64, limit int) ([]models.Earthquake, error) {
	query := `
		SELECT id, magnitude, place, time, longitude, latitude, depth 
		FROM earthquakes 
		WHERE magnitude >= $1 
		ORDER BY time DESC 
		LIMIT $2`

	rows, err := r.DB.Query(query, minMag, limit)
	if err != nil {
		return nil, fmt.Errorf("error al consultar sismos: %w", err)
	}
	defer rows.Close()

	var earthquakes []models.Earthquake
	for rows.Next() {
		var eq models.Earthquake
		var lon, lat, depth float64
		err := rows.Scan(&eq.ID, &eq.Info.Mag, &eq.Info.Place, &eq.Info.Time, &lon, &lat, &depth)
		if err != nil {
			return nil, fmt.Errorf("error al escanear sismo: %w", err)
		}
		eq.Geometry.Coordinates = []float64{lon, lat, depth}
		earthquakes = append(earthquakes, eq)
	}

	return earthquakes, nil
}

// SaveEarthquake inserta o actualiza un sismo (Upsert)
func (r *EarthquakeRepository) SaveEarthquake(eq models.Earthquake) error {
	query := `
		INSERT INTO earthquakes (id, magnitude, place, time, longitude, latitude, depth)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			magnitude = EXCLUDED.magnitude,
			place = EXCLUDED.place,
			time = EXCLUDED.time,
			longitude = EXCLUDED.longitude,
			latitude = EXCLUDED.latitude,
			depth = EXCLUDED.depth`

	lon := eq.Geometry.Coordinates[0]
	lat := eq.Geometry.Coordinates[1]
	depth := eq.Geometry.Coordinates[2]

	_, err := r.DB.Exec(query, eq.ID, eq.Info.Mag, eq.Info.Place, eq.Info.Time, lon, lat, depth)
	return err
}
