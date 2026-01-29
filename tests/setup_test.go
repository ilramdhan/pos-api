package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ilramdhan/pos-api/internal/config"
	"github.com/ilramdhan/pos-api/internal/handler"
	"github.com/ilramdhan/pos-api/internal/middleware"
	"github.com/ilramdhan/pos-api/internal/models"
	"github.com/ilramdhan/pos-api/internal/repository"
	"github.com/ilramdhan/pos-api/internal/service"
	"github.com/ilramdhan/pos-api/internal/utils"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

// TestEnv holds all test dependencies
type TestEnv struct {
	Config  *config.Config
	DB      *sql.DB
	Engine  *gin.Engine
	JWT     *utils.JWTManager
	Cookies []*http.Cookie

	// Services
	AuthService        *service.AuthService
	UserService        *service.UserService
	CategoryService    *service.CategoryService
	ProductService     *service.ProductService
	CustomerService    *service.CustomerService
	TransactionService *service.TransactionService

	// Cleanup function
	Cleanup func()
}

// SetupTestEnv creates a test environment with PostgreSQL test database
func SetupTestEnv(t *testing.T) *TestEnv {
	t.Helper()

	// Change to project root so Viper can find .env
	if err := os.Chdir(".."); err != nil {
		// Already at project root or can't change - continue anyway
	}

	// Load .env file first using Viper
	cfg := config.Load()

	// Use TEST_DB_CONN if set, otherwise fall back to DB_CONN from .env
	dbConn := os.Getenv("TEST_DB_CONN")
	if dbConn == "" {
		dbConn = cfg.Database.ConnectionString
	}

	if dbConn == "" {
		t.Fatal("No database connection string found. Set TEST_DB_CONN or DB_CONN in .env")
	}

	os.Setenv("JWT_SECRET", cfg.JWT.Secret)
	os.Setenv("APP_ENV", "test")
	os.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")

	// Connect to PostgreSQL test database
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	// Clean and setup test database
	cleanTestDatabase(t, db)
	runMigrations(t, db)
	seedTestData(t, db)

	// Setup dependencies
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	jwtManager := utils.NewJWTManager(
		cfg.JWT.Secret,
		time.Duration(cfg.JWT.ExpiryHours)*time.Hour,
		time.Duration(cfg.JWT.RefreshExpiryHours)*time.Hour,
	)

	// Repositories (PostgreSQL only)
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, jwtManager)
	userService := service.NewUserService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	productService := service.NewProductService(productRepo, categoryRepo)
	customerService := service.NewCustomerService(customerRepo)
	transactionService := service.NewTransactionService(transactionRepo, productRepo, customerRepo)

	// Setup routes
	setupTestRoutes(engine, cfg, jwtManager, authService, userService, categoryService,
		productService, customerService, transactionService, db)

	return &TestEnv{
		Config:             cfg,
		DB:                 db,
		Engine:             engine,
		JWT:                jwtManager,
		AuthService:        authService,
		UserService:        userService,
		CategoryService:    categoryService,
		ProductService:     productService,
		CustomerService:    customerService,
		TransactionService: transactionService,
		Cleanup: func() {
			cleanTestDatabase(t, db)
			db.Close()
		},
	}
}

// cleanTestDatabase drops all test data
func cleanTestDatabase(t *testing.T, db *sql.DB) {
	t.Helper()

	// Drop tables in correct order due to foreign keys
	tables := []string{
		"transaction_items",
		"transactions",
		"notifications",
		"products",
		"categories",
		"customers",
		"users",
	}

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table))
		if err != nil {
			t.Logf("Warning: failed to drop table %s: %v", table, err)
		}
	}
}

// runMigrations executes the PostgreSQL migration
func runMigrations(t *testing.T, db *sql.DB) {
	t.Helper()

	migration := `
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			name TEXT NOT NULL,
			phone TEXT DEFAULT '',
			role TEXT NOT NULL DEFAULT 'cashier' CHECK (role IN ('admin', 'manager', 'cashier')),
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS categories (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT DEFAULT '',
			slug TEXT NOT NULL UNIQUE,
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS products (
			id TEXT PRIMARY KEY,
			category_id TEXT NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
			sku TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			description TEXT DEFAULT '',
			price DECIMAL(10, 2) NOT NULL DEFAULT 0,
			stock INTEGER NOT NULL DEFAULT 0,
			image_url TEXT DEFAULT '',
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS customers (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT DEFAULT '',
			phone TEXT DEFAULT '',
			address TEXT DEFAULT '',
			loyalty_points INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS transactions (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
			customer_id TEXT REFERENCES customers(id) ON DELETE SET NULL,
			invoice_number TEXT NOT NULL UNIQUE,
			subtotal DECIMAL(12, 2) NOT NULL DEFAULT 0,
			tax_amount DECIMAL(12, 2) DEFAULT 0,
			discount_amount DECIMAL(12, 2) DEFAULT 0,
			total_amount DECIMAL(12, 2) NOT NULL DEFAULT 0,
			payment_method TEXT NOT NULL CHECK (payment_method IN ('cash', 'card', 'qris', 'transfer')),
			status TEXT NOT NULL DEFAULT 'completed' CHECK (status IN ('pending', 'completed', 'cancelled', 'refunded')),
			notes TEXT DEFAULT '',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS transaction_items (
			id TEXT PRIMARY KEY,
			transaction_id TEXT NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
			product_id TEXT NOT NULL,
			product_name TEXT NOT NULL,
			unit_price DECIMAL(10, 2) NOT NULL,
			quantity INTEGER NOT NULL,
			subtotal DECIMAL(12, 2) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS notifications (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			type TEXT NOT NULL,
			title TEXT NOT NULL,
			message TEXT NOT NULL,
			is_read BOOLEAN DEFAULT FALSE,
			action_url TEXT DEFAULT '',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := db.Exec(migration)
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
}

// setupTestRoutes configures all routes for testing
func setupTestRoutes(engine *gin.Engine, cfg *config.Config, jwtManager *utils.JWTManager,
	authService *service.AuthService, userService *service.UserService,
	categoryService *service.CategoryService, productService *service.ProductService,
	customerService *service.CustomerService, transactionService *service.TransactionService,
	db *sql.DB) {

	// Handlers
	healthHandler := handler.NewHealthHandler(cfg)
	authHandler := handler.NewAuthHandler(authService, cfg)
	userHandler := handler.NewUserHandler(userService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)
	customerHandler := handler.NewCustomerHandler(customerService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	posHandler := handler.NewPOSHandler(productService, transactionService)

	// Health
	engine.GET("/health", healthHandler.Check)

	// API v1
	v1 := engine.Group("/api/v1")
	{
		// Auth (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(jwtManager))
		{
			protected.GET("/auth/me", authHandler.Me)
			protected.PUT("/auth/me", authHandler.UpdateProfile)

			// Users (Admin only)
			users := protected.Group("/users")
			users.Use(middleware.RequireRole(models.RoleAdmin))
			{
				users.GET("", userHandler.List)
				users.GET("/:id", userHandler.Get)
				users.POST("", userHandler.Create)
				users.PUT("/:id", userHandler.Update)
				users.DELETE("/:id", userHandler.Delete)
			}

			// Categories
			categories := protected.Group("/categories")
			{
				categories.GET("", categoryHandler.List)
				categories.GET("/:id", categoryHandler.Get)
				categories.POST("", middleware.RequireRole(models.RoleAdmin, models.RoleManager), categoryHandler.Create)
				categories.PUT("/:id", middleware.RequireRole(models.RoleAdmin, models.RoleManager), categoryHandler.Update)
				categories.DELETE("/:id", middleware.RequireRole(models.RoleAdmin), categoryHandler.Delete)
			}

			// Products
			products := protected.Group("/products")
			{
				products.GET("", productHandler.List)
				products.GET("/:id", productHandler.Get)
				products.POST("", middleware.RequireRole(models.RoleAdmin, models.RoleManager), productHandler.Create)
				products.PUT("/:id", middleware.RequireRole(models.RoleAdmin, models.RoleManager), productHandler.Update)
				products.DELETE("/:id", middleware.RequireRole(models.RoleAdmin), productHandler.Delete)
				products.PATCH("/:id/stock", productHandler.UpdateStock)
			}

			// Customers
			customers := protected.Group("/customers")
			{
				customers.GET("", customerHandler.List)
				customers.GET("/:id", customerHandler.Get)
				customers.POST("", customerHandler.Create)
				customers.PUT("/:id", customerHandler.Update)
				customers.DELETE("/:id", middleware.RequireRole(models.RoleAdmin, models.RoleManager), customerHandler.Delete)
			}

			// Transactions
			transactions := protected.Group("/transactions")
			{
				transactions.GET("", transactionHandler.List)
				transactions.GET("/:id", transactionHandler.Get)
				transactions.POST("", transactionHandler.Create)
				transactions.PATCH("/:id/status", middleware.RequireRole(models.RoleAdmin, models.RoleManager), transactionHandler.UpdateStatus)
			}

			// POS
			pos := protected.Group("/pos")
			{
				pos.GET("/products", posHandler.GetProducts)
				pos.POST("/checkout", posHandler.CreateTransaction)
				pos.GET("/held", posHandler.GetHeldTransactions)
				pos.POST("/hold", posHandler.HoldTransactionCreate)
				pos.DELETE("/held/:id", posHandler.DeleteHeldTransaction)
			}
		}
	}
}

// Test data IDs (valid UUIDs)
var (
	TestAdminID    = "a0000000-0000-0000-0000-000000000001"
	TestManagerID  = "a0000000-0000-0000-0000-000000000002"
	TestCashierID  = "a0000000-0000-0000-0000-000000000003"
	TestCategoryID = "c0000000-0000-0000-0000-000000000001"
	TestProductID  = "p0000000-0000-0000-0000-000000000001"
	TestCustomerID = "d0000000-0000-0000-0000-000000000001"
)

// seedTestData adds initial test data to the database
func seedTestData(t *testing.T, db *sql.DB) {
	t.Helper()
	now := time.Now()

	// Create test users
	users := []struct {
		id       string
		email    string
		password string
		name     string
		role     string
	}{
		{TestAdminID, "admin@test.local", "Admin123!", "Test Admin", "admin"},
		{TestManagerID, "manager@test.local", "Manager123!", "Test Manager", "manager"},
		{TestCashierID, "cashier@test.local", "Cashier123!", "Test Cashier", "cashier"},
	}

	for _, u := range users {
		hash, _ := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
		_, err := db.Exec(`
			INSERT INTO users (id, email, password_hash, name, role, is_active, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, TRUE, $6, $7)
		`, u.id, u.email, string(hash), u.name, u.role, now, now)
		if err != nil {
			t.Fatalf("Failed to seed user %s: %v", u.email, err)
		}
	}

	// Create test category
	_, err := db.Exec(`
		INSERT INTO categories (id, name, description, slug, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, TRUE, $5, $6)
	`, TestCategoryID, "Test Category", "Test category description", "test-category", now, now)
	if err != nil {
		t.Fatalf("Failed to seed category: %v", err)
	}

	// Create test product
	_, err = db.Exec(`
		INSERT INTO products (id, category_id, sku, name, description, price, stock, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, TRUE, $8, $9)
	`, TestProductID, TestCategoryID, "TEST-001", "Test Product", "Test product description", 10000, 100, now, now)
	if err != nil {
		t.Fatalf("Failed to seed product: %v", err)
	}

	// Create test customer
	_, err = db.Exec(`
		INSERT INTO customers (id, name, email, phone, address, loyalty_points, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, TestCustomerID, "Test Customer", "customer@test.local", "081234567890", "Test Address", 100, now, now)
	if err != nil {
		t.Fatalf("Failed to seed customer: %v", err)
	}
}

// LoginAs performs login and returns auth cookies
func (env *TestEnv) LoginAs(t *testing.T, email, password string) []*http.Cookie {
	t.Helper()

	body := map[string]string{
		"email":    email,
		"password": password,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	env.Engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Login failed with status %d: %s", w.Code, w.Body.String())
	}

	cookies := w.Result().Cookies()
	env.Cookies = cookies
	return cookies
}

// LoginAsAdmin logs in as admin user
func (env *TestEnv) LoginAsAdmin(t *testing.T) []*http.Cookie {
	return env.LoginAs(t, "admin@test.local", "Admin123!")
}

// LoginAsManager logs in as manager user
func (env *TestEnv) LoginAsManager(t *testing.T) []*http.Cookie {
	return env.LoginAs(t, "manager@test.local", "Manager123!")
}

// LoginAsCashier logs in as cashier user
func (env *TestEnv) LoginAsCashier(t *testing.T) []*http.Cookie {
	return env.LoginAs(t, "cashier@test.local", "Cashier123!")
}

// MakeRequest makes an HTTP request with optional auth cookies
func (env *TestEnv) MakeRequest(t *testing.T, method, path string, body interface{}, cookies []*http.Cookie) *httptest.ResponseRecorder {
	t.Helper()

	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	w := httptest.NewRecorder()
	env.Engine.ServeHTTP(w, req)

	return w
}

// AssertStatus checks if response has expected status code
func AssertStatus(t *testing.T, w *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if w.Code != expected {
		t.Errorf("Expected status %d, got %d. Body: %s", expected, w.Code, w.Body.String())
	}
}

// ParseResponse parses JSON response into map
func ParseResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	return response
}

// GenerateUUID generates a new UUID for testing
func GenerateUUID() string {
	return uuid.New().String()
}
