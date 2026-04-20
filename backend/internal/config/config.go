// Package config handles application configuration from environment variables.
package config

import "os"

// Config holds the application configuration loaded from environment variables.
type Config struct {
	Port            string
	GinMode         string
	CORSAllowOrigin string
	DatabaseURL     string // NUEVO: URL de conexión a PostgreSQL
	JWTSecret       string // Clave para firmar tokens
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		Port:            getEnv("PORT", "8080"),
		GinMode:         getEnv("GIN_MODE", "debug"),
		CORSAllowOrigin: getEnv("CORS_ALLOW_ORIGIN", "*"),
		// Añadimos la lectura de la URL de base de datos.
		// El fallback asume una base local estándar sin SSL, útil para desarrollo.
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/seismic_db?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "super-secret-key-change-me"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
