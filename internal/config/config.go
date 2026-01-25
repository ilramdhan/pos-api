package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all application configuration
type Config struct {
	App      AppConfig
	JWT      JWTConfig
	Database DatabaseConfig
	RateLimit RateLimitConfig
	CORS     CORSConfig
	Log      LogConfig
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Env     string
	Port    string
	Name    string
	Version string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret             string
	ExpiryHours        time.Duration
	RefreshExpiryHours time.Duration
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver string
	Path   string
	// D1 specific
	D1DatabaseID string
	D1AccountID  string
	D1APIToken   string
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RPS   int
	Burst int
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string
	Format string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		App: AppConfig{
			Env:     getEnv("APP_ENV", "development"),
			Port:    getEnv("APP_PORT", "8080"),
			Name:    getEnv("APP_NAME", "POS API"),
			Version: getEnv("APP_VERSION", "1.0.0"),
		},
		JWT: JWTConfig{
			Secret:             getEnv("JWT_SECRET", "default-secret-change-me"),
			ExpiryHours:        time.Duration(getEnvInt("JWT_EXPIRY_HOURS", 24)) * time.Hour,
			RefreshExpiryHours: time.Duration(getEnvInt("JWT_REFRESH_EXPIRY_HOURS", 168)) * time.Hour,
		},
		Database: DatabaseConfig{
			Driver:       getEnv("DB_DRIVER", "sqlite3"),
			Path:         getEnv("DB_PATH", "./data/pos.db"),
			D1DatabaseID: getEnv("D1_DATABASE_ID", ""),
			D1AccountID:  getEnv("D1_ACCOUNT_ID", ""),
			D1APIToken:   getEnv("D1_API_TOKEN", ""),
		},
		RateLimit: RateLimitConfig{
			RPS:   getEnvInt("RATE_LIMIT_RPS", 100),
			Burst: getEnvInt("RATE_LIMIT_BURST", 200),
		},
		CORS: CORSConfig{
			AllowedOrigins: strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"), ","),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "debug"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets an environment variable as an integer with a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// IsDevelopment returns true if the app is running in development mode
func (c *Config) IsDevelopment() bool {
	return c.App.Env == "development"
}

// IsProduction returns true if the app is running in production mode
func (c *Config) IsProduction() bool {
	return c.App.Env == "production"
}
