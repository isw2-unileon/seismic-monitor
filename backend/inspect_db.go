package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		fmt.Println("DATABASE_URL is not set")
		os.Exit(1)
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		fmt.Printf("Error opening db: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// 1. Remove UNIQUE constraint from username to allow duplicate names (Full Name)
	_, err = db.Exec("ALTER TABLE users DROP CONSTRAINT IF EXISTS users_username_key;")
	if err != nil {
		fmt.Printf("Error dropping unique constraint: %v\n", err)
	} else {
		fmt.Println("Unique constraint on username dropped successfully")
	}

	// 2. Ensure min_magnitude_alert column exists in user_locations
	_, err = db.Exec("ALTER TABLE user_locations ADD COLUMN IF NOT EXISTS min_magnitude_alert NUMERIC(3,1) NOT NULL DEFAULT 3.0;")
	if err != nil {
		fmt.Printf("Error adding column to user_locations: %v\n", err)
	} else {
		fmt.Println("Column min_magnitude_alert ensured in user_locations")
	}

	// 3. Inspect columns of user_locations
	rows, err := db.Query("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = 'user_locations';")
	if err == nil {
		fmt.Println("\nColumns in user_locations:")
		for rows.Next() {
			var name, dtype string
			rows.Scan(&name, &dtype)
			fmt.Printf("- %s (%s)\n", name, dtype)
		}
		rows.Close()
	}
}
