package repository

import (
	"context"
	"database/sql"

	"github.com/ilramdhan/pos-api/internal/models"
	"github.com/ilramdhan/pos-api/internal/utils"
)

type categoryRepository struct {
	db *sql.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *models.Category) error {
	query := `
		INSERT INTO categories (id, name, description, slug, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		category.ID, category.Name, category.Description, category.Slug,
		category.IsActive, category.CreatedAt, category.UpdatedAt,
	)
	return err
}

func (r *categoryRepository) GetByID(ctx context.Context, id string) (*models.Category, error) {
	query := `
		SELECT id, name, description, slug, is_active, created_at, updated_at
		FROM categories WHERE id = ?
	`
	category := &models.Category{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&category.ID, &category.Name, &category.Description, &category.Slug,
		&category.IsActive, &category.CreatedAt, &category.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return category, err
}

func (r *categoryRepository) GetBySlug(ctx context.Context, slug string) (*models.Category, error) {
	query := `
		SELECT id, name, description, slug, is_active, created_at, updated_at
		FROM categories WHERE slug = ?
	`
	category := &models.Category{}
	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&category.ID, &category.Name, &category.Description, &category.Slug,
		&category.IsActive, &category.CreatedAt, &category.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return category, err
}

func (r *categoryRepository) Update(ctx context.Context, category *models.Category) error {
	query := `
		UPDATE categories SET name = ?, description = ?, slug = ?, is_active = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		category.Name, category.Description, category.Slug, category.IsActive,
		category.UpdatedAt, category.ID,
	)
	return err
}

func (r *categoryRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM categories WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *categoryRepository) List(ctx context.Context, pagination utils.Pagination) ([]*models.Category, int, error) {
	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM categories`
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := `
		SELECT id, name, description, slug, is_active, created_at, updated_at
		FROM categories
		ORDER BY ` + pagination.OrderBy() + `
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, pagination.Limit(), pagination.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var categories []*models.Category
	for rows.Next() {
		category := &models.Category{}
		if err := rows.Scan(
			&category.ID, &category.Name, &category.Description, &category.Slug,
			&category.IsActive, &category.CreatedAt, &category.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		categories = append(categories, category)
	}

	return categories, total, rows.Err()
}
