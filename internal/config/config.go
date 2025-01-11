package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Database    DatabaseConfig
	Server      ServerConfig
	Logger      LoggerConfig
	Features    FeatureFlags
	Backup      BackupConfig
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

type BackupConfig struct {
	Enabled    bool
	Directory  string
	Interval   time.Duration
	MaxBackups int
}

func LoadEnv() error {
	return godotenv.Load()
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
		Backup: BackupConfig{
			Enabled:    getEnvAsBool("BACKUP_ENABLED", false),
			Directory:  getEnv("BACKUP_DIRECTORY", "./backups"),
			Interval:   getEnvAsDuration("BACKUP_INTERVAL", 24*time.Hour),
			MaxBackups: getEnvAsInt("BACKUP_MAX_BACKUPS", 30),
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

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}
