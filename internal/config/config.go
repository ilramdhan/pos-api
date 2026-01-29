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

	// Get database connection string with fallback support for Zeabur
	dbConn := getDBConnectionString()

	return &Config{
		App: AppConfig{
			Env:     viper.GetString("APP_ENV"),
			Port:    getPort(),
			Name:    viper.GetString("APP_NAME"),
			Version: viper.GetString("APP_VERSION"),
		},
		JWT: JWTConfig{
			Secret:             viper.GetString("JWT_SECRET"),
			ExpiryHours:        viper.GetInt("JWT_EXPIRY_HOURS"),
			RefreshExpiryHours: viper.GetInt("JWT_REFRESH_EXPIRY_HOURS"),
		},
		Database: DatabaseConfig{
			ConnectionString: dbConn,
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

// getDBConnectionString gets the database connection string with fallback support
// Priority: DB_CONN > POSTGRES_URI > POSTGRES_CONNECTION_STRING > construct from components
// getDBConnectionString gets the database connection string with fallback support
// Priority: DB_CONN > POSTGRES_URI > POSTGRES_CONNECTION_STRING > construct from components
func getDBConnectionString() string {
	var conn string

	// 1. Try environment variables in order
	if c := viper.GetString("DB_CONN"); c != "" {
		conn = c
	} else if c := viper.GetString("POSTGRES_URI"); c != "" {
		conn = c
	} else if c := viper.GetString("POSTGRES_CONNECTION_STRING"); c != "" {
		conn = c
	}

	if conn != "" {
		return conn
	}

	// 2. Fallback: Construct from Zeabur/standard components
	host := viper.GetString("POSTGRES_HOST")
	if host == "" {
		host = viper.GetString("POSTGRESQL_HOST")
	}

	if host != "" {
		port := viper.GetString("POSTGRES_PORT")
		if port == "" {
			port = "5432"
		}
		user := viper.GetString("POSTGRES_USERNAME")
		if user == "" {
			user = "postgres"
		}
		password := viper.GetString("POSTGRES_PASSWORD")
		database := viper.GetString("POSTGRES_DATABASE")
		if database == "" {
			database = "postgres"
		}

		// Default to sslmode=require for security
		return "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + database + "?sslmode=require"
	}

	log.Println("WARNING: No database connection string found")
	return ""
}

// getPort gets port with fallback support for Zeabur
func getPort() string {
	// Zeabur uses PORT
	if port := viper.GetString("PORT"); port != "" {
		return port
	}
	return viper.GetString("APP_PORT")
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
