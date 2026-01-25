package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ilramdhan/pos-api/internal/dto"
	"github.com/ilramdhan/pos-api/internal/models"
	"github.com/ilramdhan/pos-api/internal/repository"
	"github.com/ilramdhan/pos-api/internal/utils"
)

// ProductService handles product operations
type ProductService struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
}

// NewProductService creates a new product service
func NewProductService(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

// Create creates a new product
func (s *ProductService) Create(ctx context.Context, req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	// Validate category exists
	category, err := s.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	// Check SKU uniqueness
	existing, err := s.productRepo.GetBySKU(ctx, req.SKU)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("SKU already exists")
	}

	now := time.Now()
	product := &models.Product{
		ID:          uuid.New().String(),
		CategoryID:  req.CategoryID,
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	product.Category = category
	return s.toResponse(product), nil
}

// GetByID retrieves a product by ID
func (s *ProductService) GetByID(ctx context.Context, id string) (*dto.ProductResponse, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	return s.toResponse(product), nil
}

// Update updates a product
func (s *ProductService) Update(ctx context.Context, id string, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	// Validate category exists if provided
	var category *models.Category
	if req.CategoryID != "" {
		category, err = s.categoryRepo.GetByID(ctx, req.CategoryID)
		if err != nil {
			return nil, err
		}
		if category == nil {
			return nil, errors.New("category not found")
		}
		product.CategoryID = req.CategoryID
	}

	// Check SKU uniqueness if changed
	if req.SKU != "" && product.SKU != req.SKU {
		existing, err := s.productRepo.GetBySKU(ctx, req.SKU)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, errors.New("SKU already exists")
		}
		product.SKU = req.SKU
	}

	// Update only fields that are provided (partial update)
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price != nil && *req.Price > 0 {
		product.Price = *req.Price
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.ImageURL != "" {
		product.ImageURL = req.ImageURL
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}
	product.UpdatedAt = time.Now()

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	if category != nil {
		product.Category = category
	}
	return s.toResponse(product), nil
}

// Delete deletes a product
func (s *ProductService) Delete(ctx context.Context, id string) error {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("product not found")
	}

	return s.productRepo.Delete(ctx, id)
}

// UpdateStock updates product stock
func (s *ProductService) UpdateStock(ctx context.Context, id string, req *dto.UpdateStockRequest) (*dto.ProductResponse, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	var newStock int
	switch req.Operation {
	case "add":
		newStock = product.Stock + req.Quantity
	case "subtract":
		newStock = product.Stock - req.Quantity
		if newStock < 0 {
			return nil, errors.New("insufficient stock")
		}
	case "set":
		newStock = req.Quantity
	default:
		return nil, errors.New("invalid operation")
	}

	if err := s.productRepo.UpdateStock(ctx, id, newStock); err != nil {
		return nil, err
	}

	product.Stock = newStock
	return s.toResponse(product), nil
}

// List lists products with pagination and filters
func (s *ProductService) List(ctx context.Context, filter dto.ProductListFilter, pagination utils.Pagination) ([]*dto.ProductResponse, int, error) {
	products, total, err := s.productRepo.List(ctx, filter, pagination)
	if err != nil {
		return nil, 0, err
	}

	var responses []*dto.ProductResponse
	for _, product := range products {
		responses = append(responses, s.toResponse(product))
	}

	return responses, total, nil
}

func (s *ProductService) toResponse(product *models.Product) *dto.ProductResponse {
	resp := &dto.ProductResponse{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		SKU:         product.SKU,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		ImageURL:    product.ImageURL,
		IsActive:    product.IsActive,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}

	if product.Category != nil {
		resp.Category = &dto.CategoryResponse{
			ID:          product.Category.ID,
			Name:        product.Category.Name,
			Description: product.Category.Description,
			Slug:        product.Category.Slug,
			IsActive:    product.Category.IsActive,
			CreatedAt:   product.Category.CreatedAt,
			UpdatedAt:   product.Category.UpdatedAt,
		}
	}

	return resp
}
