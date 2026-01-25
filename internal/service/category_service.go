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

// CategoryService handles category operations
type CategoryService struct {
	categoryRepo repository.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(categoryRepo repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

// Create creates a new category
func (s *CategoryService) Create(ctx context.Context, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	// Check slug uniqueness
	existing, err := s.categoryRepo.GetBySlug(ctx, req.Slug)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("slug already exists")
	}

	now := time.Now()
	category := &models.Category{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Slug:        req.Slug,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return s.toResponse(category), nil
}

// GetByID retrieves a category by ID
func (s *CategoryService) GetByID(ctx context.Context, id string) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	return s.toResponse(category), nil
}

// Update updates a category
func (s *CategoryService) Update(ctx context.Context, id string, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	// Check slug uniqueness if changed
	if req.Slug != "" && category.Slug != req.Slug {
		existing, err := s.categoryRepo.GetBySlug(ctx, req.Slug)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, errors.New("slug already exists")
		}
		category.Slug = req.Slug
	}

	// Update only fields that are provided (partial update)
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}
	category.UpdatedAt = time.Now()

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	return s.toResponse(category), nil
}

// Delete deletes a category
func (s *CategoryService) Delete(ctx context.Context, id string) error {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if category == nil {
		return errors.New("category not found")
	}

	return s.categoryRepo.Delete(ctx, id)
}

// List lists categories with pagination
func (s *CategoryService) List(ctx context.Context, pagination utils.Pagination) ([]*dto.CategoryResponse, int, error) {
	categories, total, err := s.categoryRepo.List(ctx, pagination)
	if err != nil {
		return nil, 0, err
	}

	var responses []*dto.CategoryResponse
	for _, category := range categories {
		responses = append(responses, s.toResponse(category))
	}

	return responses, total, nil
}

func (s *CategoryService) toResponse(category *models.Category) *dto.CategoryResponse {
	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		Slug:        category.Slug,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}
