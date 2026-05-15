package database

import (
	"database/sql"
	"fmt"
	"time"

	"seismic-monitor/backend/internal/models"

	"golang.org/x/crypto/bcrypt"
)

// UserRepository define las operaciones de persistencia para usuarios
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository crea una nueva instancia de UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// CreateUser inserta un nuevo usuario en la base de datos
func (r *UserRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (username, email, password_hash, location, alert_radius_km, min_magnitude_alert, created_at) VALUES ($1, $2, $3, ST_SetSRID(ST_MakePoint($4, $5), 4326), $6, $7, $8) RETURNING id`

	user.CreatedAt = time.Now()
	username := user.Email // Usar el email como username ya que username es requerido

	err := r.DB.QueryRow(query, username, user.Email, user.PasswordHash, user.Longitude, user.Latitude, user.AlertRadius, user.MinMagnitude, user.CreatedAt).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("error al crear usuario: %w", err)
	}
	return nil
}

// FindUserByEmail busca un usuario por su dirección de correo electrónico
func (r *UserRepository) FindUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, password_hash, ST_Y(location::geometry) as latitude, ST_X(location::geometry) as longitude, alert_radius_km, min_magnitude_alert, created_at FROM users WHERE email = $1`
	err := r.DB.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Latitude, &user.Longitude, &user.AlertRadius, &user.MinMagnitude, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil // Usuario no encontrado
	} else if err != nil {
		return nil, fmt.Errorf("error al buscar usuario por email: %w", err)
	}
	return user, nil
}

// UpdateUserLocation actualiza la latitud, longitud y radio de alerta de un usuario
func (r *UserRepository) UpdateUserLocation(userID string, latitude, longitude, alertRadius, minMagnitude float64) error {
	query := `UPDATE users SET location = ST_SetSRID(ST_MakePoint($1, $2), 4326), alert_radius_km = $3, min_magnitude_alert = $4 WHERE id = $5`
	err := r.DB.QueryRow(query, longitude, latitude, alertRadius, minMagnitude, userID).Err()
	if err != nil {
		return fmt.Errorf("error al actualizar la ubicación del usuario: %w", err)
	}
	return nil
}

// HashPassword genera un hash bcrypt de la contraseña
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compara una contraseña en texto plano con su hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetAffectedUsers implementa la interfaz SpatialRepository usando PostGIS.
// Busca usuarios cuyo radio de alerta cubra la ubicación del sismo.
func (r *UserRepository) GetAffectedUsers(sismo models.Feature) ([]models.User, error) {
	// Query que usa ST_DWithin para calcular si el sismo está dentro del radio del usuario.
	// ST_DWithin usa metros con 'geography', así que multiplicamos km * 1000.
	query := `
		SELECT id, email, alert_radius_km, min_magnitude_alert
		FROM users
		WHERE ST_DWithin(
			location::geography, 
			ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, 
			alert_radius_km * 1000
		) AND min_magnitude_alert <= $3`

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
		// Nota: solo escaneamos los campos que devuelve la query (id, email, radius, magnitude)
		if err := rows.Scan(&u.ID, &u.Email, &u.AlertRadius, &u.MinMagnitude); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

// GetUsersNearLocation busca usuarios que tengan un punto geográfico dentro de su radio de alerta
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
