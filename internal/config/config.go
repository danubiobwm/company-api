package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppPort  string
	DB       DBConfig
	LogLevel string
	Env      string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func Load() Config {
	return Config{
		AppPort:  getenv("APP_PORT", "8080"),
		LogLevel: getenv("LOG_LEVEL", "info"),
		Env:      getenv("ENV", "development"),
		DB: DBConfig{
			Host:     getenv("DATABASE_HOST", "db"),
			Port:     getenv("DATABASE_PORT", "5432"),
			User:     getenv("DATABASE_USER", "postgres"),
			Password: getenv("DATABASE_PASSWORD", "postgres"),
			DBName:   getenv("DATABASE_NAME", "companydb"),
			SSLMode:  getenv("DATABASE_SSLMODE", "disable"),
		},
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getenvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}
