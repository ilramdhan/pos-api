package repository

import (
	"context"

	"github.com/ilramdhan/pos-api/internal/dto"
	"github.com/ilramdhan/pos-api/internal/models"
	"github.com/ilramdhan/pos-api/internal/utils"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, role string, pagination utils.Pagination) ([]*models.User, int, error)
}

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
	Create(ctx context.Context, category *models.Category) error
	GetByID(ctx context.Context, id string) (*models.Category, error)
	GetBySlug(ctx context.Context, slug string) (*models.Category, error)
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, pagination utils.Pagination) ([]*models.Category, int, error)
}

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	GetByID(ctx context.Context, id string) (*models.Product, error)
	GetBySKU(ctx context.Context, sku string) (*models.Product, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id string) error
	UpdateStock(ctx context.Context, id string, quantity int) error
	List(ctx context.Context, filter dto.ProductListFilter, pagination utils.Pagination) ([]*models.Product, int, error)
}

// CustomerRepository defines the interface for customer data access
type CustomerRepository interface {
	Create(ctx context.Context, customer *models.Customer) error
	GetByID(ctx context.Context, id string) (*models.Customer, error)
	Update(ctx context.Context, customer *models.Customer) error
	Delete(ctx context.Context, id string) error
	UpdateLoyaltyPoints(ctx context.Context, id string, points int) error
	List(ctx context.Context, filter dto.CustomerListFilter, pagination utils.Pagination) ([]*models.Customer, int, error)
}

// TransactionRepository defines the interface for transaction data access
type TransactionRepository interface {
	Create(ctx context.Context, transaction *models.Transaction) error
	GetByID(ctx context.Context, id string) (*models.Transaction, error)
	GetByInvoiceNumber(ctx context.Context, invoiceNumber string) (*models.Transaction, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	List(ctx context.Context, filter dto.TransactionListFilter, pagination utils.Pagination) ([]*models.Transaction, int, error)
	GetDailySales(ctx context.Context, dateFrom, dateTo string) ([]dto.DailySalesReport, error)
	GetMonthlySales(ctx context.Context, dateFrom, dateTo string) ([]dto.MonthlySalesReport, error)
	GetTopProducts(ctx context.Context, limit int, dateFrom, dateTo string) ([]dto.TopProductReport, error)
}

// TransactionItemRepository defines the interface for transaction item data access
type TransactionItemRepository interface {
	Create(ctx context.Context, item *models.TransactionItem) error
	GetByTransactionID(ctx context.Context, transactionID string) ([]*models.TransactionItem, error)
}
