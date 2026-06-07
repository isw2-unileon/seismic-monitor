package models

import (
	"time"
)

type User struct {
	ID           string        `json:"id"`
	Email        string        `json:"email" binding:"required,email"`
	Name         string        `json:"name"`
	PasswordHash string        `json:"-"` // Hide in JSON
	Latitude     float64       `json:"latitude"`
	Longitude    float64       `json:"longitude"`
	AlertRadius  float64       `json:"alert_radius"`
	MinMagnitude float64       `json:"min_magnitude"`
	AlertCenters []AlertCenter `json:"alert_centers,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
}

type AlertCenter struct {
	ID           string  `json:"id"`
	UserID       string  `json:"user_id"`
	Latitude     float64 `json:"lat"`
	Longitude    float64 `json:"lng"`
	Radius       float64 `json:"radius"`
	MinMagnitude float64 `json:"min_magnitude"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateLocationRequest struct {
	Name         string   `json:"name"`
	Latitude     *float64 `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude    *float64 `json:"longitude" binding:"required,min=-180,max=180"`
	AlertRadius  *float64 `json:"alert_radius" binding:"required,min=1"` // Radius in kilometers
	MinMagnitude *float64 `json:"min_magnitude" binding:"required,min=0,max=10"`
}


