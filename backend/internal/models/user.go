package models

import (
	"time"
)

// User representa la estructura de un usuario en el sistema
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email" binding:"required,email"`
	PasswordHash string    `json:"-"` // Ocultar en JSON
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	AlertRadius  float64   `json:"alert_radius"`
	CreatedAt    time.Time `json:"created_at"`
}

// RegisterRequest es el DTO para el registro de usuarios
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest es el DTO para el login de usuarios
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateLocationRequest es el DTO para actualizar la ubicación del usuario
type UpdateLocationRequest struct {
	Latitude    float64 `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude   float64 `json:"longitude" binding:"required,min=-180,max=180"`
	AlertRadius float64 `json:"alert_radius" binding:"required,min=1"` // Radio en kilómetros
}
