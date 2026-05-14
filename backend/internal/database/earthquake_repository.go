package database

import (
	"database/sql"
	"fmt"
	"seismic-monitor/backend/internal/models"
	"time"
)

type EarthquakeRepository struct {
	DB *sql.DB
}

func NewEarthquakeRepository(db *sql.DB) *EarthquakeRepository {
	return &EarthquakeRepository{DB: db}
}

// GetFilteredEarthquakes obtiene sismos filtrados por magnitud mínima y límite
func (r *EarthquakeRepository) GetFilteredEarthquakes(minMag float64, limit int) ([]models.Feature, error) {
	query := `
		SELECT usgs_id, richter_scale, place_name, ocurred_at, ST_X(location::geometry), ST_Y(location::geometry), depth_km 
		FROM earthquake 
		WHERE richter_scale >= $1 
		ORDER BY ocurred_at DESC 
		LIMIT $2`

	rows, err := r.DB.Query(query, minMag, limit)
	if err != nil {
		return nil, fmt.Errorf("error al consultar sismos: %w", err)
	}
	defer rows.Close()

	var earthquakes []models.Feature
	for rows.Next() {
		var eq models.Feature
		var lon, lat, depth float64
		var ocurredAt time.Time
		err := rows.Scan(&eq.ID, &eq.Info.Mag, &eq.Info.Place, &ocurredAt, &lon, &lat, &depth)
		if err != nil {
			return nil, fmt.Errorf("error al escanear sismo: %w", err)
		}
		eq.Info.Time = ocurredAt.UnixNano() / int64(time.Millisecond)
		eq.Geometry.Coordinates = []float64{lon, lat, depth}
		earthquakes = append(earthquakes, eq)
	}

	return earthquakes, nil
}

// SaveEarthquake inserta o actualiza un sismo (Upsert)
func (r *EarthquakeRepository) SaveEarthquake(eq models.Feature) error {
	query := `
		INSERT INTO earthquake (usgs_id, richter_scale, place_name, ocurred_at, location, depth_km)
		VALUES ($1, $2, $3, $4, ST_SetSRID(ST_MakePoint($5, $6), 4326), $7)
		ON CONFLICT (usgs_id) DO UPDATE SET
			richter_scale = EXCLUDED.richter_scale,
			place_name = EXCLUDED.place_name,
			ocurred_at = EXCLUDED.ocurred_at,
			location = EXCLUDED.location,
			depth_km = EXCLUDED.depth_km`

	lon := eq.Geometry.Coordinates[0]
	lat := eq.Geometry.Coordinates[1]
	depth := eq.Geometry.Coordinates[2]

	// Convert Unix ms to time.Time
	t := time.Unix(0, eq.Info.Time*int64(time.Millisecond))

	_, err := r.DB.Exec(query, eq.ID, eq.Info.Mag, eq.Info.Place, t, lon, lat, depth)
	return err
}

func (r *EarthquakeRepository) GetEarthquakesSince(since time.Time) ([]models.Feature, error) {
	query := `
		SELECT usgs_id, richter_scale, place_name, ocurred_at, ST_X(location::geometry), ST_Y(location::geometry), depth_km
		FROM earthquake
		WHERE ocurred_at >= $1
		ORDER BY ocurred_at DESC`

	rows, err := r.DB.Query(query, since)
	if err != nil {
		return nil, fmt.Errorf("error al consultar sismos recientes: %w", err)
	}
	defer rows.Close()

	var earthquakes []models.Feature
	for rows.Next() {
		var eq models.Feature
		var lon, lat, depth float64
		var ocurredAt time.Time
		err := rows.Scan(&eq.ID, &eq.Info.Mag, &eq.Info.Place, &ocurredAt, &lon, &lat, &depth)
		if err != nil {
			return nil, fmt.Errorf("error al escanear sismo: %w", err)
		}
		eq.Info.Time = ocurredAt.UnixNano() / int64(time.Millisecond)
		eq.Geometry.Coordinates = []float64{lon, lat, depth}
		earthquakes = append(earthquakes, eq)
	}

	return earthquakes, nil
}
