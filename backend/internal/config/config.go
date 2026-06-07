package config

import (
	"log"
	"os"

	"github.com/joho/godotenv" // Necessary to read the .env file
)

type Config struct {
	DatabaseURL string
	Port        string
	GinMode     string
	JWTSecret   string
}

func Load() *Config {
	// THIS IS THE MISSING MAGIC LINE!
	// godotenv.Load() finds the .env file and loads it into the system variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: No se encontró archivo .env, leyendo variables del sistema...")
	}

	return &Config{
		// Now getEnv will find the variable loaded from .env
		DatabaseURL: getEnv("DATABASE_URL", "postgres://db:db@localhost:5432/db?sslmode=disable"),
		Port:        getEnv("API_PORT", "8081"),
		GinMode:     getEnv("GIN_MODE", "debug"),
		JWTSecret:   getEnv("JWT_SECRET", "tu_secreto_super_seguro_por_defecto"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
