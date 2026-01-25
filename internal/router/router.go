package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ilramdhan/pos-api/internal/config"
	"github.com/ilramdhan/pos-api/internal/database"
	"github.com/ilramdhan/pos-api/internal/handler"
	"github.com/ilramdhan/pos-api/internal/middleware"
	"github.com/ilramdhan/pos-api/internal/models"
	"github.com/ilramdhan/pos-api/internal/repository"
	"github.com/ilramdhan/pos-api/internal/service"
	"github.com/ilramdhan/pos-api/internal/utils"
)

// Router holds all route handlers
type Router struct {
	Engine *gin.Engine
	cfg    *config.Config
}

// New creates and configures a new router
func New(cfg *config.Config, db *database.Database) *Router {
	// Set Gin mode
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()

	// Global middleware
	engine.Use(gin.Recovery())
	engine.Use(middleware.LoggerMiddleware())
	engine.Use(middleware.RequestIDMiddleware())
	engine.Use(middleware.CORSMiddleware(cfg.CORS.AllowedOrigins))

	// Rate limiter
	rateLimiter := middleware.NewIPRateLimiter(cfg.RateLimit.RPS, cfg.RateLimit.Burst)
	rateLimiter.StartCleanup(5*time.Minute, 10*time.Minute)
	engine.Use(middleware.RateLimitMiddleware(rateLimiter))

	// JWT Manager
	jwtManager := utils.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpiryHours, cfg.JWT.RefreshExpiryHours)

	// Repositories
	userRepo := repository.NewUserRepository(db.DB)
	categoryRepo := repository.NewCategoryRepository(db.DB)
	productRepo := repository.NewProductRepository(db.DB)
	customerRepo := repository.NewCustomerRepository(db.DB)
	transactionRepo := repository.NewTransactionRepository(db.DB)

	// Services
	authService := service.NewAuthService(userRepo, jwtManager)
	categoryService := service.NewCategoryService(categoryRepo)
	productService := service.NewProductService(productRepo, categoryRepo)
	customerService := service.NewCustomerService(customerRepo)
	transactionService := service.NewTransactionService(transactionRepo, productRepo, customerRepo)
	reportService := service.NewReportService(transactionRepo)

	// Handlers
	healthHandler := handler.NewHealthHandler(cfg)
	authHandler := handler.NewAuthHandler(authService, cfg)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)
	customerHandler := handler.NewCustomerHandler(customerService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	reportHandler := handler.NewReportHandler(reportService)
	dashboardHandler := handler.NewDashboardHandler(
		transactionService,
		productService,
		categoryService,
		customerService,
		reportService,
	)

	// Routes
	// Health check (public)
	engine.GET("/health", healthHandler.Check)

	// API v1
	v1 := engine.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", healthHandler.CheckV1)

		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(jwtManager))
		{
			// Auth (protected)
			protected.GET("/auth/me", authHandler.Me)
			protected.PUT("/auth/me", authHandler.UpdateProfile)
			protected.GET("/auth/me/activity", authHandler.GetActivityLog)

			// Categories
			categories := protected.Group("/categories")
			{
				categories.GET("", categoryHandler.List)
				categories.GET("/stats", dashboardHandler.GetCategoryStats)
				categories.GET("/activity-log", dashboardHandler.GetCategoryActivityLog)
				categories.GET("/:id", categoryHandler.Get)
				categories.POST("", middleware.RequireRole(models.RoleAdmin, models.RoleManager), categoryHandler.Create)
				categories.PUT("/:id", middleware.RequireRole(models.RoleAdmin, models.RoleManager), categoryHandler.Update)
				categories.DELETE("/:id", middleware.RequireRole(models.RoleAdmin), categoryHandler.Delete)
			}

			// Products
			products := protected.Group("/products")
			{
				products.GET("", productHandler.List)
				products.GET("/stats", dashboardHandler.GetProductStats)
				products.GET("/stock-movements", dashboardHandler.GetStockMovements)
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
				customers.GET("/stats", dashboardHandler.GetCustomerStats)
				customers.GET("/:id", customerHandler.Get)
				customers.POST("", customerHandler.Create)
				customers.PUT("/:id", customerHandler.Update)
				customers.DELETE("/:id", middleware.RequireRole(models.RoleAdmin, models.RoleManager), customerHandler.Delete)
			}

			// Transactions
			transactions := protected.Group("/transactions")
			{
				transactions.GET("", transactionHandler.List)
				transactions.GET("/recent", dashboardHandler.GetRecentTransactions)
				transactions.GET("/:id", transactionHandler.Get)
				transactions.POST("", transactionHandler.Create)
				transactions.PATCH("/:id/status", middleware.RequireRole(models.RoleAdmin, models.RoleManager), transactionHandler.UpdateStatus)
			}

			// Reports
			reports := protected.Group("/reports")
			{
				// Dashboard (all authenticated users)
				reports.GET("/dashboard/stats", dashboardHandler.GetDashboardStats)
				reports.GET("/sales/realtime", dashboardHandler.GetRealtimeSales)

				// Detailed reports (admin/manager only)
				reports.GET("/sales/daily", middleware.RequireRole(models.RoleAdmin, models.RoleManager), reportHandler.DailySales)
				reports.GET("/sales/monthly", middleware.RequireRole(models.RoleAdmin, models.RoleManager), reportHandler.MonthlySales)
				reports.GET("/products/top", middleware.RequireRole(models.RoleAdmin, models.RoleManager), reportHandler.TopProducts)
			}

			// System
			system := protected.Group("/system")
			{
				system.GET("/health/detailed", dashboardHandler.GetSystemHealth)
			}
		}
	}

	return &Router{
		Engine: engine,
		cfg:    cfg,
	}
}

// Run starts the HTTP server
func (r *Router) Run() error {
	return r.Engine.Run(":" + r.cfg.App.Port)
}
