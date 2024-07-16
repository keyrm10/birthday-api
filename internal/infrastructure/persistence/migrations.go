package persistence

import (
	"database/sql"
	"log"
)

func RunMigrations(db *sql.DB) error {
	log.Println("Running database migrations...")

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			username VARCHAR(255) PRIMARY KEY,
			date_of_birth DATE NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_users_date_of_birth ON users(date_of_birth);
	`)
	if err != nil {
		log.Printf("Error creating users table or index: %v", err)
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}
