package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL   string
	ServerAddress string
}

func Load() (*Config, error) {
	pgUser := os.Getenv("PGUSER")
	pgPassword := os.Getenv("PGPASSWORD")
	pgHost := os.Getenv("PGHOST")
	pgPort := os.Getenv("PGPORT")
	pgDatabase := os.Getenv("PGDATABASE")
	pgSSLMode := os.Getenv("PGSSLMODE")

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		pgUser, pgPassword, pgHost, pgPort, pgDatabase, pgSSLMode)

	return &Config{
		DatabaseURL:   databaseURL,
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
	}, nil
}
