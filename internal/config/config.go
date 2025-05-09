package config

import (
    "fmt"
    "log"
    "os"
    "strconv"
    "github.com/joho/godotenv"
)

type Config struct {
    Database DatabaseConfig
    Auth0    Auth0Config    `yaml:"auth0"`
}

type DatabaseConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    User     string `json:"user"`
    Password string `json:"password"`
    DBName   string `json:"dbname"`
}

// Auth0Config holds Auth0-specific configuration
type Auth0Config struct {
    Domain    string `yaml:"domain"`
    Audience  string `yaml:"audience"`
    Issuer    string `yaml:"issuer"`
    Algorithms []string `yaml:"algorithms"`
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

    // Load Auth0 configuration
    domain := getEnv("AUTH0_DOMAIN", "")
    issuer := getEnv("AUTH0_ISSUER", "")
    if issuer == "" && domain != "" {
        issuer = fmt.Sprintf("https://%s/", domain)
    }

    return &Config{
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     port,
            User:     getEnv("DB_USER", ""),
            Password: getEnv("DB_PASSWORD", ""),
            DBName:   getEnv("DB_NAME", ""),
        },
        Auth0: Auth0Config{
            Domain:    domain,
            Audience: getEnv("AUTH0_AUDIENCE", ""),
            Issuer:   issuer,
            Algorithms: []string{"RS256"},
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