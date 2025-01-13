package config

import (
    "log"
    "os"
    "github.com/joho/godotenv"
)


type Config struct {
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

func Load() (*Config, error) {
    // Add this at the start of the function
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: Error loading .env file: %v", err)
    }

    // Add debug logging
    log.Printf("DB_HOST: %s", os.Getenv("DB_HOST"))
    log.Printf("DB_USER: %s", os.Getenv("DB_USER"))
    log.Printf("DB_NAME: %s", os.Getenv("DB_NAME"))

    // Rest of your existing code...
    return &Config{
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     parseInt(getEnv("DB_PORT", "5432")),
            User:     getEnv("DB_USER", "postgres"),
            Password: getEnv("DB_PASSWORD", "postgres"),
            DBName:   getEnv("DB_NAME", "horse_tracking"),
        },
    }, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func parseInt(s string) int {
	var i int
	for _, c := range s {
		i = i*10 + int(c-'0')
	}
	return i
}
