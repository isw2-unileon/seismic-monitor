package config

import (
	"os"
)

// Config almacena toda la configuración de la aplicación
type Config struct {
	DatabaseURL string
	Port        string
	GinMode     string
	JWTSecret   string
}

// Load lee las variables de entorno o usa valores por defecto
func Load() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://db:db@localhost:5432/db?sslmode=disable"),
		Port:        getEnv("PORT", "8081"),
		GinMode:     getEnv("GIN_MODE", "debug"), // Por defecto en modo debug para desarrollo local
		JWTSecret:   getEnv("JWT_SECRET", "tu_secreto_super_seguro_por_defecto"),
	}
}

// Función auxiliar para leer variables con un valor por defecto
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
