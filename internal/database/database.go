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

	// Add phone column if it doesn't exist (SQLite doesn't support IF NOT EXISTS for columns)
	_, err = d.Exec("SELECT phone FROM users LIMIT 1")
	if err != nil {
		// Column doesn't exist, add it
		_, err = d.Exec("ALTER TABLE users ADD COLUMN phone TEXT DEFAULT ''")
		if err != nil {
			// Ignore error if column already exists
			fmt.Printf("Note: phone column may already exist: %v\n", err)
		}
	}

	// Create notifications table
	_, err = d.Exec(`
		CREATE TABLE IF NOT EXISTS notifications (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			type TEXT NOT NULL,
			title TEXT NOT NULL,
			message TEXT NOT NULL,
			is_read INTEGER DEFAULT 0,
			action_url TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		fmt.Printf("Note: notifications table may already exist: %v\n", err)
	}

	// Create indexes for notifications
	d.Exec("CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id)")
	d.Exec("CREATE INDEX IF NOT EXISTS idx_notifications_unread ON notifications(user_id, is_read)")

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
