package models

import "time"

type UserReport struct {
	ID         string    `json:"reported_earthquake_id"`
	Longitude  float64   `json:"longitude" binding:"required"`
	Latitude   float64   `json:"latitude" binding:"required"`
	ReportedAt time.Time `json:"reported_at"`
}
