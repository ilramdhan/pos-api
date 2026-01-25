package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ilramdhan/pos-api/internal/dto"
	"github.com/ilramdhan/pos-api/internal/models"
	"github.com/ilramdhan/pos-api/internal/repository"
	"github.com/ilramdhan/pos-api/internal/utils"
)

const (
	// TaxRate is the default tax rate (10%)
	TaxRate = 0.10
	// LoyaltyPointsPerTransaction is points earned per transaction
	LoyaltyPointsPerTransaction = 10
)

// TransactionService handles transaction operations
type TransactionService struct {
	transactionRepo repository.TransactionRepository
	productRepo     repository.ProductRepository
	customerRepo    repository.CustomerRepository
}

// NewTransactionService creates a new transaction service
func NewTransactionService(
	transactionRepo repository.TransactionRepository,
	productRepo repository.ProductRepository,
	customerRepo repository.CustomerRepository,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		productRepo:     productRepo,
		customerRepo:    customerRepo,
	}
}

// Create creates a new transaction (sale)
func (s *TransactionService) Create(ctx context.Context, userID string, req *dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
	now := time.Now()
	transactionID := uuid.New().String()
	invoiceNumber := fmt.Sprintf("INV-%s-%s", now.Format("20060102"), transactionID[:8])

	// Build transaction items and calculate totals
	var items []models.TransactionItem
	var subtotal float64

	for _, itemReq := range req.Items {
		product, err := s.productRepo.GetByID(ctx, itemReq.ProductID)
		if err != nil {
			return nil, err
		}
		if product == nil {
			return nil, fmt.Errorf("product %s not found", itemReq.ProductID)
		}
		if !product.IsActive {
			return nil, fmt.Errorf("product %s is not available", product.Name)
		}
		if !product.HasSufficientStock(itemReq.Quantity) {
			return nil, fmt.Errorf("insufficient stock for product %s", product.Name)
		}

		itemSubtotal := product.Price * float64(itemReq.Quantity)
		item := models.TransactionItem{
			ID:            uuid.New().String(),
			TransactionID: transactionID,
			ProductID:     product.ID,
			ProductName:   product.Name,
			UnitPrice:     product.Price,
			Quantity:      itemReq.Quantity,
			Subtotal:      itemSubtotal,
			CreatedAt:     now,
		}
		items = append(items, item)
		subtotal += itemSubtotal

		// Update product stock
		newStock := product.Stock - itemReq.Quantity
		if err := s.productRepo.UpdateStock(ctx, product.ID, newStock); err != nil {
			return nil, err
		}
	}

	taxAmount := subtotal * TaxRate
	totalAmount := subtotal + taxAmount - req.DiscountAmount

	transaction := &models.Transaction{
		ID:             transactionID,
		UserID:         userID,
		CustomerID:     req.CustomerID,
		InvoiceNumber:  invoiceNumber,
		Subtotal:       subtotal,
		TaxAmount:      taxAmount,
		DiscountAmount: req.DiscountAmount,
		TotalAmount:    totalAmount,
		PaymentMethod:  req.PaymentMethod,
		Status:         models.StatusCompleted,
		Notes:          req.Notes,
		CreatedAt:      now,
		UpdatedAt:      now,
		Items:          items,
	}

	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, err
	}

	// Add loyalty points to customer if specified
	if req.CustomerID != nil && *req.CustomerID != "" {
		customer, err := s.customerRepo.GetByID(ctx, *req.CustomerID)
		if err == nil && customer != nil {
			newPoints := customer.LoyaltyPoints + LoyaltyPointsPerTransaction
			_ = s.customerRepo.UpdateLoyaltyPoints(ctx, *req.CustomerID, newPoints)
		}
	}

	return s.toResponse(transaction), nil
}

// GetByID retrieves a transaction by ID
func (s *TransactionService) GetByID(ctx context.Context, id string) (*dto.TransactionResponse, error) {
	transaction, err := s.transactionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if transaction == nil {
		return nil, errors.New("transaction not found")
	}

	return s.toResponse(transaction), nil
}

// UpdateStatus updates a transaction's status
func (s *TransactionService) UpdateStatus(ctx context.Context, id string, req *dto.UpdateTransactionStatusRequest) (*dto.TransactionResponse, error) {
	transaction, err := s.transactionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if transaction == nil {
		return nil, errors.New("transaction not found")
	}

	// Validate status transition
	switch req.Status {
	case models.StatusCancelled:
		if !transaction.IsCancellable() {
			return nil, errors.New("transaction cannot be cancelled")
		}
		// Restore stock for cancelled transactions
		for _, item := range transaction.Items {
			product, err := s.productRepo.GetByID(ctx, item.ProductID)
			if err == nil && product != nil {
				newStock := product.Stock + item.Quantity
				_ = s.productRepo.UpdateStock(ctx, item.ProductID, newStock)
			}
		}
	case models.StatusRefunded:
		if !transaction.IsRefundable() {
			return nil, errors.New("transaction cannot be refunded")
		}
		// Restore stock for refunded transactions
		for _, item := range transaction.Items {
			product, err := s.productRepo.GetByID(ctx, item.ProductID)
			if err == nil && product != nil {
				newStock := product.Stock + item.Quantity
				_ = s.productRepo.UpdateStock(ctx, item.ProductID, newStock)
			}
		}
	}

	if err := s.transactionRepo.UpdateStatus(ctx, id, req.Status); err != nil {
		return nil, err
	}

	transaction.Status = req.Status
	return s.toResponse(transaction), nil
}

// List lists transactions with pagination and filters
func (s *TransactionService) List(ctx context.Context, filter dto.TransactionListFilter, pagination utils.Pagination) ([]*dto.TransactionResponse, int, error) {
	transactions, total, err := s.transactionRepo.List(ctx, filter, pagination)
	if err != nil {
		return nil, 0, err
	}

	var responses []*dto.TransactionResponse
	for _, transaction := range transactions {
		responses = append(responses, s.toResponse(transaction))
	}

	return responses, total, nil
}

func (s *TransactionService) toResponse(transaction *models.Transaction) *dto.TransactionResponse {
	resp := &dto.TransactionResponse{
		ID:             transaction.ID,
		UserID:         transaction.UserID,
		CustomerID:     transaction.CustomerID,
		InvoiceNumber:  transaction.InvoiceNumber,
		Subtotal:       transaction.Subtotal,
		TaxAmount:      transaction.TaxAmount,
		DiscountAmount: transaction.DiscountAmount,
		TotalAmount:    transaction.TotalAmount,
		PaymentMethod:  transaction.PaymentMethod,
		Status:         transaction.Status,
		Notes:          transaction.Notes,
		CreatedAt:      transaction.CreatedAt,
		UpdatedAt:      transaction.UpdatedAt,
	}

	if transaction.User != nil {
		resp.User = &dto.UserResponse{
			ID:       transaction.User.ID,
			Email:    transaction.User.Email,
			Name:     transaction.User.Name,
			Role:     transaction.User.Role,
			IsActive: transaction.User.IsActive,
		}
	}

	for _, item := range transaction.Items {
		resp.Items = append(resp.Items, dto.TransactionItemResponse{
			ID:          item.ID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			UnitPrice:   item.UnitPrice,
			Quantity:    item.Quantity,
			Subtotal:    item.Subtotal,
		})
	}

	return resp
}
