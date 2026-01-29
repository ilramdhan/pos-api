package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ilramdhan/pos-api/internal/config"
	_ "github.com/lib/pq"
)

// Database wraps the sql.DB connection for PostgreSQL (Supabase)
type Database struct {
	*sql.DB
}

// New creates a new PostgreSQL/Supabase database connection
func New(cfg *config.DatabaseConfig) (*Database, error) {
	if cfg.ConnectionString == "" {
		return nil, fmt.Errorf("DB_CONN is required for PostgreSQL/Supabase connection")
	}

	// Force IPv4 resolution to fix Zeabur "network unreachable" (IPv6) issues
	resolvedConnStr, err := resolveToIPv4(cfg.ConnectionString)
	if err != nil {
		log.Printf("⚠️ Failed to resolve hostname to IPv4: %v. Using original connection string.", err)
		resolvedConnStr = cfg.ConnectionString
	} else {
		log.Println("✓ Resolved hostname to IPv4 to bypass potential IPv6 routing issues")
	}

	log.Println("Connecting to PostgreSQL/Supabase database...")
	// Log connection string with password redacted for debugging
	redactedConn := redactPassword(resolvedConnStr)
	log.Printf("Connection string: %s", redactedConn)

	db, err := sql.Open("postgres", resolvedConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open PostgreSQL database: %w", err)
	}

	// PostgreSQL connection pool settings optimized for Supabase
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	// Test connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✓ PostgreSQL/Supabase database connected successfully")
	return &Database{DB: db}, nil
}

// resolveToIPv4 parses the connection string, resolves the hostname to an IPv4 address,
// and returns a new connection string with the IP address.
// This is necessary because some environments (like Zeabur) have broken IPv6 routing to Supabase.
func resolveToIPv4(connStr string) (string, error) {
	// Simple parsing to find hostname
	// Format: postgres://user:pass@hostname:port/db...
	// We want to replace hostname with IPv4
	
// resolveToIPv4 parses the connection string, resolves the hostname to an IPv4 address,
// and returns a new connection string with the IP address.
// This is necessary because some environments (like Zeabur) have broken IPv6 routing to Supabase.
func resolveToIPv4(connStr string) (string, error) {
	u, err := url.Parse(connStr)
	if err != nil {
		return "", err
	}

	host, port, _ := net.SplitHostPort(u.Host)
	if host == "" {
		host = u.Host // No port specified
	}

	// If it's already an IP, pass through
	if net.ParseIP(host) != nil {
		return connStr, nil
	}

	// Resolve IPs
	ips, err := net.LookupIP(host)
	if err != nil {
		return "", err
	}

	var ipv4 net.IP
	for _, ip := range ips {
		if ip.To4() != nil {
			ipv4 = ip
			break
		}
	}

	if ipv4 == nil {
		return "", fmt.Errorf("no IPv4 address found for %s", host)
	}

	// Reconstruct URL with IPv4
	if port != "" {
		u.Host = ipv4.String() + ":" + port
	} else {
		u.Host = ipv4.String()
	}

	return u.String(), nil
}

func redactPassword(connStr string) string {
	u, err := url.Parse(connStr)
	if err != nil {
		return "invalid-url"
	}
	if _, hasPassword := u.User.Password(); hasPassword {
		u.User = url.UserPassword(u.User.Username(), "REDACTED")
	}
	return u.String()
}

// Migrate runs database migrations from PostgreSQL migration file
func (d *Database) Migrate(migrationsPath string) error {
	migrationFile := filepath.Join(migrationsPath, "001_init_postgres.sql")

	content, err := os.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Execute migration
	if _, err := d.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	log.Println("✓ Database migrations applied successfully")
	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.DB.Close()
}

// Transaction executes a function within a database transaction
func (d *Database) Transaction(fn func(tx *sql.Tx) error) error {
	tx, err := d.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback error: %v, original error: %w", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
