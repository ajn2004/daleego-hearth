package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port              string
	DatabaseURL       string
	AdminAPIKey       string
	PairingCodeSecret string
}

func Load() (Config, error) {

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbIP := os.Getenv("DB_IP")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")
	AdminAPIKey := os.Getenv("ADMIN_API_KEY")
	PairingCodeSecret := os.Getenv("PAIRING_CODE_SECRET")

	if dbUser == "" {
		return Config{}, fmt.Errorf("DB_USER environment variable is required")
	}

	if dbPass == "" {
		return Config{}, fmt.Errorf("DB_PASS environment variable is required")
	}

	if dbIP == "" {
		return Config{}, fmt.Errorf("DB_IP environment variable is required")
	}

	if dbPort == "" {
		return Config{}, fmt.Errorf("DB_PORT environment variable is required")
	}

	if dbName == "" {
		return Config{}, fmt.Errorf("DB_NAME environment variable is required")
	}
	if dbSSLMode == "" {
		dbSSLMode = "disable"
	}
	if AdminAPIKey == "" {
		return Config{}, fmt.Errorf("ADMIN_API_KEY environment variable is required")
	}
	if PairingCodeSecret == "" {
		return Config{}, fmt.Errorf("PAIRING_CODE_SECRET environment variable is required")
	}

	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser,
		dbPass,
		dbIP,
		dbPort,
		dbName,
		dbSSLMode,
	)
	cfg := Config{
		Port:              os.Getenv("PORT"),
		DatabaseURL:       databaseURL,
		AdminAPIKey:       AdminAPIKey,
		PairingCodeSecret: PairingCodeSecret,
	}

	if cfg.Port == "" {
		cfg.Port = "8080" // Default port
	}

	return cfg, nil
}
