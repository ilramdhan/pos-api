package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ilramdhan/pos-api/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

// Database wraps the sql.DB connection
type Database struct {
	*sql.DB
}

// New creates a new database connection
func New(cfg *config.DatabaseConfig) (*Database, error) {
	// Ensure data directory exists
	if cfg.Driver == "sqlite3" {
		dir := filepath.Dir(cfg.Path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create data directory: %w", err)
		}
	}

	db, err := sql.Open(cfg.Driver, cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys for SQLite
	if cfg.Driver == "sqlite3" {
		if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
			return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
		}
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{DB: db}, nil
}

// Migrate runs database migrations
func (d *Database) Migrate(migrationsPath string) error {
	// Read migration file
	migrationFile := filepath.Join(migrationsPath, "001_init.sql")
	content, err := os.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Execute migration
	if _, err := d.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

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
