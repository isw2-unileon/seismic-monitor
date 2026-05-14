package config

import (
	"log"
	"os"

	"github.com/joho/godotenv" // Necesario para leer el archivo .env
)

type Config struct {
	DatabaseURL string
	Port        string
	GinMode     string
	JWTSecret   string
}

func Load() *Config {
	// ¡ESTA ES LA LÍNEA MÁGICA QUE FALTA!
	// godotenv.Load() busca el archivo .env y lo carga en las variables del sistema
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: No se encontró archivo .env, leyendo variables del sistema...")
	}

	return &Config{
		// Ahora getEnv sí encontrará la variable cargada del .env
		DatabaseURL: getEnv("DATABASE_URL", "postgres://db:db@localhost:5432/db?sslmode=disable"),
		Port:        getEnv("API_PORT", "8081"),
		GinMode:     getEnv("GIN_MODE", "debug"),
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
