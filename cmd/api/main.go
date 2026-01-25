package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ilramdhan/pos-api/internal/config"
	"github.com/ilramdhan/pos-api/internal/database"
	"github.com/ilramdhan/pos-api/internal/router"
	"github.com/joho/godotenv"
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
