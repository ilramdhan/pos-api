package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	App       AppConfig
	JWT       JWTConfig
	Database  DatabaseConfig
	RateLimit RateLimitConfig
	CORS      CORSConfig
}

// AppConfig holds application-level configuration
type AppConfig struct {
	Env     string
	Port    string
	Name    string
	Version string
}

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	Secret             string
	ExpiryHours        int
	RefreshExpiryHours int
}

// DatabaseConfig holds PostgreSQL/Supabase database configuration
type DatabaseConfig struct {
	ConnectionString string
}

// RateLimitConfig holds rate limiter configuration
type RateLimitConfig struct {
	RPS   int
	Burst int
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
}

// Load loads configuration from environment variables using Viper
func Load() *Config {
	// Set up Viper
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// Read .env file (optional - environment variables take precedence)
	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Set defaults
	setDefaults()

	return &Config{
		App: AppConfig{
			Env:     viper.GetString("APP_ENV"),
			Port:    viper.GetString("APP_PORT"),
			Name:    viper.GetString("APP_NAME"),
			Version: viper.GetString("APP_VERSION"),
		},
		JWT: JWTConfig{
			Secret:             viper.GetString("JWT_SECRET"),
			ExpiryHours:        viper.GetInt("JWT_EXPIRY_HOURS"),
			RefreshExpiryHours: viper.GetInt("JWT_REFRESH_EXPIRY_HOURS"),
		},
		Database: DatabaseConfig{
			ConnectionString: viper.GetString("DB_CONN"),
		},
		RateLimit: RateLimitConfig{
			RPS:   viper.GetInt("RATE_LIMIT_RPS"),
			Burst: viper.GetInt("RATE_LIMIT_BURST"),
		},
		CORS: CORSConfig{
			AllowedOrigins: parseOrigins(viper.GetString("CORS_ALLOWED_ORIGINS")),
		},
	}
}

// setDefaults sets default configuration values
func setDefaults() {
	// Application defaults
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("APP_NAME", "GoPOS API")
	viper.SetDefault("APP_VERSION", "2.0.0")

	// JWT defaults
	viper.SetDefault("JWT_SECRET", "change-this-secret-in-production")
	viper.SetDefault("JWT_EXPIRY_HOURS", 24)
	viper.SetDefault("JWT_REFRESH_EXPIRY_HOURS", 168) // 7 days

	// Rate limiter defaults
	viper.SetDefault("RATE_LIMIT_RPS", 100)
	viper.SetDefault("RATE_LIMIT_BURST", 200)

	// CORS defaults
	viper.SetDefault("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173")
}

// parseOrigins parses comma-separated origins string into slice
func parseOrigins(origins string) []string {
	if origins == "" {
		return []string{"http://localhost:3000"}
	}
	parts := strings.Split(origins, ",")
	result := make([]string, 0, len(parts))
	for _, origin := range parts {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			result = append(result, origin)
		}
	}
	return result
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return c.App.Env == "production"
}

// IsDevelopment returns true if running in development environment
func (c *Config) IsDevelopment() bool {
	return c.App.Env == "development"
}
