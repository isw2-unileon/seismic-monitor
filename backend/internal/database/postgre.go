package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // Driver de Postgres
)

// Connect establece y configura el pool de conexiones
func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("error al abrir la base de datos: %w", err)
	}

	// Configuración de buenas prácticas para el pool
	db.SetMaxOpenConns(25)                 // Máximo de conexiones abiertas
	db.SetMaxIdleConns(25)                 // Máximo de conexiones inactivas
	db.SetConnMaxLifetime(5 * time.Minute) // Tiempo de vida de la conexión

	// Verificar si la conexión es válida
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error de ping a la base de datos: %w", err)
	}

	return db, nil
}
