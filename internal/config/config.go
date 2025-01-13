package config

import (
    "log"
    "os"
    "strconv"
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
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: Error loading .env file: %v", err)
    }

    // Debug logging
    log.Printf("DB_HOST: %s", os.Getenv("DB_HOST"))
    log.Printf("DB_USER: %s", os.Getenv("DB_USER"))
    log.Printf("DB_NAME: %s", os.Getenv("DB_NAME"))

    port, err := parseInt(getEnv("DB_PORT", "5432"))
    if err != nil {
        log.Printf("Warning: Invalid DB_PORT value, using default: %v", err)
        port = 5432
    }

    return &Config{
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     port,
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

func parseInt(s string) (int, error) {
    return strconv.Atoi(s)
}