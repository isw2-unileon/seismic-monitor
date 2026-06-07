package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"seismic-monitor/backend/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (username, email, password_hash, location, alert_radius_km, min_magnitude_alert, created_at) VALUES ($1, $2, $3, ST_SetSRID(ST_MakePoint($4, $5), 4326)::geography, $6, $7, $8) RETURNING id`

	user.CreatedAt = time.Now()
	// If the name is empty, use the email as the default username
	username := user.Name
	if username == "" {
		username = user.Email
	}

	err := r.DB.QueryRow(query, username, user.Email, user.PasswordHash, user.Longitude, user.Latitude, user.AlertRadius, user.MinMagnitude, user.CreatedAt).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("error al crear usuario: %w", err)
	}
	return nil
}

func (r *UserRepository) FindUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, username, password_hash, ST_Y(location::geometry) as latitude, ST_X(location::geometry) as longitude, alert_radius_km, min_magnitude_alert, created_at FROM users WHERE email = $1`
	err := r.DB.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.Latitude, &user.Longitude, &user.AlertRadius, &user.MinMagnitude, &user.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil // User not found
	} else if err != nil {
		return nil, fmt.Errorf("error al buscar usuario por email: %w", err)
	}

	centers, err := r.GetUserAlertCenters(user.ID)
	if err == nil {
		user.AlertCenters = centers
	}

	return user, nil
}

func (r *UserRepository) GetUserAlertCenters(userID string) ([]models.AlertCenter, error) {
	query := `SELECT id, ST_Y(location::geometry) as lat, ST_X(location::geometry) as lng, alert_radius_km, min_magnitude_alert FROM user_locations WHERE user_id = $1::uuid`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener centros de alerta: %w", err)
	}
	defer rows.Close()

	var centers []models.AlertCenter
	for rows.Next() {
		var c models.AlertCenter
		c.UserID = userID
		if err := rows.Scan(&c.ID, &c.Latitude, &c.Longitude, &c.Radius, &c.MinMagnitude); err != nil {
			return nil, err
		}
		centers = append(centers, c)
	}
	return centers, nil
}

func (r *UserRepository) DeleteUserAlertCenter(userID, centerID string) error {
	query := `DELETE FROM user_locations WHERE id = $1::uuid AND user_id = $2::uuid`
	res, err := r.DB.Exec(query, centerID, userID)
	if err != nil {
		slog.Error("Database error in DeleteUserAlertCenter", "error", err, "userID", userID, "centerID", centerID)
		return fmt.Errorf("error al eliminar centro de alerta de la BD: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no se encontró la ubicación con ID %s para el usuario %s", centerID, userID)
	}
	slog.Info("Alert center deleted from DB", "id", centerID, "userID", userID)
	return nil
}

func (r *UserRepository) AddUserAlertCenter(userID string, latitude, longitude, radius, minMagnitude float64) (string, error) {
	var newID string
	query := `INSERT INTO user_locations (user_id, location, alert_radius_km, min_magnitude_alert) VALUES ($1::uuid, ST_SetSRID(ST_MakePoint($2, $3), 4326)::geography, $4, $5) RETURNING id`
	err := r.DB.QueryRow(query, userID, longitude, latitude, radius, minMagnitude).Scan(&newID)
	if err != nil {
		slog.Error("Database error in AddUserAlertCenter", "error", err, "userID", userID, "lat", latitude, "lng", longitude)
		return "", fmt.Errorf("error al añadir centro de alerta: %w", err)
	}
	slog.Info("Alert center added to DB", "id", newID, "userID", userID)
	return newID, nil
}

func (r *UserRepository) UpdateUserLocation(userID string, name string, latitude, longitude, alertRadius, minMagnitude float64) error {
	query := `UPDATE users SET username = $1, location = ST_SetSRID(ST_MakePoint($2, $3), 4326)::geography, alert_radius_km = $4, min_magnitude_alert = $5 WHERE id = $6::uuid`
	res, err := r.DB.Exec(query, name, longitude, latitude, alertRadius, minMagnitude, userID)
	if err != nil {
		slog.Error("Database error in UpdateUserLocation", "error", err, "userID", userID, "name", name)
		return fmt.Errorf("error al actualizar el perfil del usuario: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no se encontró al usuario con ID %s", userID)
	}
	slog.Info("User profile updated in DB", "userID", userID, "rowsAffected", rows)
	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Find users whose alert radius covers the earthquake location (checks users and user_locations tables).
func (r *UserRepository) GetAffectedUsers(sismo models.Feature) ([]models.User, error) {
	// Query using ST_DWithin to calculate if the earthquake is within the user's radius.
	// Search in both the user's main location and their additional zones.
	// Each zone has its own minimum magnitude.
	query := `
		SELECT DISTINCT u.id, u.email, u.alert_radius_km, u.min_magnitude_alert
		FROM users u
		LEFT JOIN user_locations ul ON u.id = ul.user_id
		WHERE (
			(ST_DWithin(u.location::geography, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, u.alert_radius_km * 1000) AND u.min_magnitude_alert <= $3)
			OR
			(ul.location IS NOT NULL AND ST_DWithin(ul.location::geography, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, ul.alert_radius_km * 1000) AND ul.min_magnitude_alert <= $3)
		)`

	if len(sismo.Geometry.Coordinates) < 2 {
		return nil, fmt.Errorf("no se pueden calcular usuarios afectados: el sismo no tiene coordenadas válidas")
	}

	lon := sismo.Geometry.Coordinates[0]
	lat := sismo.Geometry.Coordinates[1]
	mag := sismo.Info.Mag

	rows, err := r.DB.Query(query, lon, lat, mag)
	if err != nil {
		return nil, fmt.Errorf("error buscando usuarios afectados: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		// Note: we only scan the fields returned by the query (id, email, radius, magnitude)
		if err := rows.Scan(&u.ID, &u.Email, &u.AlertRadius, &u.MinMagnitude); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepository) GetUsersNearLocation(lon, lat float64) ([]models.User, error) {
	query := `
		SELECT id, email, alert_radius_km
		FROM users
		WHERE ST_DWithin(
			location::geography, 
			ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, 
			alert_radius_km * 1000
		)`

	rows, err := r.DB.Query(query, lon, lat)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Email, &u.AlertRadius); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
