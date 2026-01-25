package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/ilramdhan/pos-api/internal/config"
	"github.com/ilramdhan/pos-api/internal/database"
	"github.com/ilramdhan/pos-api/internal/router"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load .env file (optional for development)
	_ = godotenv.Load()

	// Load configuration
	cfg := config.Load()

	log.Printf("Starting %s v%s in %s mode\n", cfg.App.Name, cfg.App.Version, cfg.App.Env)

	// Connect to database
	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate("internal/database/migrations"); err != nil {
		log.Printf("Warning: Migration failed: %v", err)
	}

	// Auto-seed if database is empty
	if shouldAutoSeed(db.DB) {
		log.Println("Database is empty, running auto-seed...")
		autoSeed(db.DB)
	}

	// Setup router
	r := router.New(cfg, db)

	// Start server in goroutine
	go func() {
		log.Printf("Server listening on port %s", cfg.App.Port)
		if err := r.Run(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}

// shouldAutoSeed checks if database needs seeding
func shouldAutoSeed(db *sql.DB) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return true // Table might not exist
	}
	return count == 0
}

// autoSeed seeds the database with demo data
func autoSeed(db *sql.DB) {
	ctx := context.Background()
	now := time.Now()

	// Seed users
	users := []struct {
		email    string
		password string
		name     string
		role     string
	}{
		{"admin@gopos.local", "Admin123!", "Admin User", "admin"},
		{"manager@gopos.local", "Manager123!", "Manager User", "manager"},
		{"cashier@gopos.local", "Cashier123!", "Cashier User", "cashier"},
	}

	for _, u := range users {
		hash, _ := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
		_, err := db.ExecContext(ctx, `
			INSERT OR IGNORE INTO users (id, email, password_hash, name, role, is_active, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, 1, ?, ?)
		`, uuid.New().String(), u.email, string(hash), u.name, u.role, now, now)
		if err == nil {
			log.Printf("  ✓ Seeded user: %s", u.email)
		}
	}

	// Seed categories
	categories := []struct {
		id   string
		name string
		slug string
	}{
		{uuid.New().String(), "Food", "food"},
		{uuid.New().String(), "Beverages", "beverages"},
		{uuid.New().String(), "Snacks", "snacks"},
		{uuid.New().String(), "Electronics", "electronics"},
	}

	catIDs := make(map[string]string)
	for _, c := range categories {
		catIDs[c.slug] = c.id
		_, err := db.ExecContext(ctx, `
			INSERT OR IGNORE INTO categories (id, name, description, slug, is_active, created_at, updated_at)
			VALUES (?, ?, ?, ?, 1, ?, ?)
		`, c.id, c.name, c.name+" items", c.slug, now, now)
		if err == nil {
			log.Printf("  ✓ Seeded category: %s", c.name)
		}
	}

	// Seed products
	products := []struct {
		catSlug string
		sku     string
		name    string
		price   float64
		stock   int
	}{
		{"food", "FOOD-001", "Nasi Goreng Spesial", 28000, 100},
		{"food", "FOOD-002", "Mie Goreng Seafood", 32000, 80},
		{"food", "FOOD-003", "Ayam Bakar Madu", 38000, 50},
		{"beverages", "BEV-001", "Kopi Hitam", 10000, 200},
		{"beverages", "BEV-002", "Kopi Susu Gula Aren", 18000, 150},
		{"beverages", "BEV-003", "Es Teh Manis", 8000, 200},
		{"snacks", "SNK-001", "Keripik Singkong", 12000, 150},
		{"snacks", "SNK-002", "Cokelat Silverqueen", 18000, 100},
		{"electronics", "ELEC-001", "USB Cable Type-C", 35000, 50},
	}

	for _, p := range products {
		catID := catIDs[p.catSlug]
		_, err := db.ExecContext(ctx, `
			INSERT OR IGNORE INTO products (id, category_id, sku, name, description, price, stock, image_url, is_active, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, '', 1, ?, ?)
		`, uuid.New().String(), catID, p.sku, p.name, p.name, p.price, p.stock, now, now)
		if err == nil {
			log.Printf("  ✓ Seeded product: %s", p.name)
		}
	}

	// Seed customers
	customers := []struct {
		name   string
		email  string
		phone  string
		points int
	}{
		{"John Doe", "john@email.com", "081234567890", 150},
		{"Jane Smith", "jane@email.com", "081234567891", 200},
		{"Bob Wilson", "bob@email.com", "081234567892", 75},
	}

	for _, c := range customers {
		_, err := db.ExecContext(ctx, `
			INSERT OR IGNORE INTO customers (id, name, email, phone, address, loyalty_points, created_at, updated_at)
			VALUES (?, ?, ?, ?, '', ?, ?, ?)
		`, uuid.New().String(), c.name, c.email, c.phone, c.points, now, now)
		if err == nil {
			log.Printf("  ✓ Seeded customer: %s", c.name)
		}
	}

	log.Println("Auto-seed completed!")
	log.Println("")
	log.Println("Test accounts:")
	log.Println("  admin@gopos.local    / Admin123!")
	log.Println("  manager@gopos.local  / Manager123!")
	log.Println("  cashier@gopos.local  / Cashier123!")
}
