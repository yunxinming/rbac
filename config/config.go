package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port      string
	JWTSecret string
	DBPath    string
}

var AppConfig *Config

func LoadConfig() {
	port := getEnv("PORT", "8080")
	jwtSecret := getEnv("JWT_SECRET", "default-secret-change-in-production")
	dbPath := getEnv("DB_PATH", "./data/rbac.db")

	AppConfig = &Config{
		Port:      port,
		JWTSecret: jwtSecret,
		DBPath:    dbPath,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
