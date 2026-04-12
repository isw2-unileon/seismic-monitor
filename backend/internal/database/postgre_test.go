package database

import (
	"os"
	"testing"
)

func TestConnectPostgres(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("Saltando test: DATABASE_URL no configurada")
	}

	db, err := Connect(dsn)
	if err != nil {
		t.Fatalf("Error al conectar a la DB: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Fatalf("No se pudo hacer ping a la base de datos: %v", err)
	}
}
