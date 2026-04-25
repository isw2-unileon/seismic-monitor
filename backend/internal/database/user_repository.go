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
	query := `INSERT INTO users (email, password_hash, latitude, longitude, alert_radius, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err := r.DB.QueryRow(query, user.Email, user.PasswordHash, user.Latitude, user.Longitude, user.AlertRadius, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("error al crear usuario: %w", err)
	}
	return nil
}

// FindUserByEmail busca un usuario por su dirección de correo electrónico
func (r *UserRepository) FindUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, password_hash, latitude, longitude, alert_radius, created_at, updated_at FROM users WHERE email = $1`
	err := r.DB.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Latitude, &user.Longitude, &user.AlertRadius, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil // Usuario no encontrado
	} else if err != nil {
		return nil, fmt.Errorf("error al buscar usuario por email: %w", err)
	}
	return user, nil
}

// UpdateUserLocation actualiza la latitud, longitud y radio de alerta de un usuario
func (r *UserRepository) UpdateUserLocation(userID int, latitude, longitude, alertRadius float64) error {
	query := `UPDATE users SET latitude = $1, longitude = $2, alert_radius = $3, updated_at = $4 WHERE id = $5`
	err := r.DB.QueryRow(query, latitude, longitude, alertRadius, time.Now(), userID).Err()
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
