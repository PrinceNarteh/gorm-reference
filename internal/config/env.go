package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var Envs = initConfig()

func initConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env file")
	}

	return &Config{
		App: appConfig{
			Env:     getEnv("APP_ENV", "development"),
			Port:    getEnvAtInt("APP_PORT", 4000),
			Version: getEnv("APP_VERSION", "1.0.0"),
		},
		DB: dbConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnvAtInt("DB_PORT", 5432),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "password"),
			Name:            getEnv("DB_NAME", "gorm"),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			MaxIdleConns:    getEnvAtInt("DB_MAX_IDLE_CONNS", 10),
			MaxOpenConns:    getEnvAtInt("DB_MAX_OPEN_CONNS", 100),
			MaxIdleTime:     getEnvAsDuration("DB_MAX_IDLE_TIME", 10*time.Minute),
			MaxConnLifetime: getEnvAsDuration("DB_MAX_CONN_LIFETIME", time.Hour),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}

func getEnvAtInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if duration, err := time.ParseDuration(valueStr); err == nil {
		return duration
	}
	return fallback
}
