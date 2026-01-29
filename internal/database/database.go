package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/ilramdhan/pos-api/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
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

	log.Println("Connecting to PostgreSQL/Supabase database using pgx driver...")
	// Log connection string with password redacted for debugging
	redactedConn := redactPassword(cfg.ConnectionString)
	log.Printf("Connection string: %s", redactedConn)

	// Use "pgx" driver instead of "postgres" (lib/pq)
	db, err := sql.Open("pgx", cfg.ConnectionString)
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
