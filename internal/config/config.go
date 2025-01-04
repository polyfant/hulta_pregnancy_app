package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Environment string
	Database    DatabaseConfig
	Server      ServerConfig
	Logger      LoggerConfig
	Features    FeatureFlags
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type ServerConfig struct {
	Port int
	Host string
	Mode string
}

type LoggerConfig struct {
	Level string
	Path  string
}

type FeatureFlags struct {
	EnableAuditLogging bool
	EnableCaching      bool
	StrictMode         bool
}

func LoadConfig() *Config {
	env := getEnv("APP_ENV", "development")
	
	config := &Config{
		Environment: env,
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "horse_tracker"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			Port: getEnvAsInt("SERVER_PORT", 8080),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Logger: LoggerConfig{
			Level: getEnv("LOG_LEVEL", "info"),
			Path:  getEnv("LOG_PATH", "./logs"),
		},
	}

	// Environment-specific feature flags
	switch env {
	case "production":
		config.Features = FeatureFlags{
			EnableAuditLogging: true,
			EnableCaching:      true,
			StrictMode:         true,
		}
	case "development":
		config.Features = FeatureFlags{
			EnableAuditLogging: false,
			EnableCaching:      false,
			StrictMode:         false,
		}
	case "testing":
		config.Features = FeatureFlags{
			EnableAuditLogging: true,
			EnableCaching:      false,
			StrictMode:         true,
		}
	}

	return config
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

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := strings.ToLower(getEnv(key, strconv.FormatBool(defaultValue)))
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// LoadTestConfig loads configuration for testing environment
func LoadTestConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "testuser",
			Password: "testpassword",
			DBName:   "horse_tracking_test",
			SSLMode:  "disable",
		},
		Server: ServerConfig{
			Host: "localhost",
			Port: 8081,
			Mode: string(gin.TestMode),
		},
		Logger: LoggerConfig{
			Path:  "./logs/test.log",
			Level: "debug",
		},
		Features: FeatureFlags{
			EnableAuditLogging: false,
			EnableCaching:      false,
			StrictMode:         false,
		},
	}
}
