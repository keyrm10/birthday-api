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
	pgUser := getEnv("PGUSER", "")
	pgPassword := getEnv("PGPASSWORD", "")
	pgHost := getEnv("PGHOST", "localhost")
	pgPort := getEnv("PGPORT", "5432")
	pgDatabase := getEnv("PGDATABASE", "")
	pgSSLMode := getEnv("PGSSLMODE", "disable")

	if pgUser == "" || pgPassword == "" || pgDatabase == "" {
		return nil, fmt.Errorf("missing required database configuration")
	}

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		pgUser, pgPassword, pgHost, pgPort, pgDatabase, pgSSLMode)

	serverAddress := getEnv("SERVER_ADDRESS", ":8080")

	return &Config{
		DatabaseURL:   databaseURL,
		ServerAddress: serverAddress,
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
