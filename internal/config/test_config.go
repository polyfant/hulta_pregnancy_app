package config

import (
	"os"
	"strconv"
)

type TestConfig struct {
	Database DatabaseConfig
}

func LoadTestConfig() *TestConfig {
	return &TestConfig{
		Database: DatabaseConfig{
			Host:     getEnv("TEST_DB_HOST", "localhost"),
			Port:     getEnvAsInt("TEST_DB_PORT", 5432),
			User:     getEnv("TEST_DB_USER", "horse_tracker_app"),
			Password: getEnv("TEST_DB_PASSWORD", "silverback_secure_password"),
			DBName:   getEnv("TEST_DB_NAME", "horse_tracking_test"),
			SSLMode:  getEnv("TEST_DB_SSLMODE", "disable"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
