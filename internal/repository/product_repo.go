package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ilramdhan/pos-api/internal/dto"
	"github.com/ilramdhan/pos-api/internal/models"
	"github.com/ilramdhan/pos-api/internal/utils"
)

type productRepository struct {
	db *sql.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
	query := `
		INSERT INTO products (id, category_id, sku, name, description, price, stock, image_url, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		product.ID, product.CategoryID, product.SKU, product.Name, product.Description,
		product.Price, product.Stock, product.ImageURL, product.IsActive,
		product.CreatedAt, product.UpdatedAt,
	)
	return err
}

func (r *productRepository) GetByID(ctx context.Context, id string) (*models.Product, error) {
	query := `
		SELECT p.id, p.category_id, p.sku, p.name, COALESCE(p.description, ''), p.price, p.stock, COALESCE(p.image_url, ''), p.is_active, p.created_at, p.updated_at,
		       c.id, c.name, COALESCE(c.description, ''), c.slug, c.is_active, c.created_at, c.updated_at
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = ?
	`
	product := &models.Product{}
	category := &models.Category{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID, &product.CategoryID, &product.SKU, &product.Name, &product.Description,
		&product.Price, &product.Stock, &product.ImageURL, &product.IsActive,
		&product.CreatedAt, &product.UpdatedAt,
		&category.ID, &category.Name, &category.Description, &category.Slug,
		&category.IsActive, &category.CreatedAt, &category.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	product.Category = category
	return product, nil
}

func (r *productRepository) GetBySKU(ctx context.Context, sku string) (*models.Product, error) {
	query := `
		SELECT id, category_id, sku, name, COALESCE(description, ''), price, stock, COALESCE(image_url, ''), is_active, created_at, updated_at
		FROM products WHERE sku = ?
	`
	product := &models.Product{}
	err := r.db.QueryRowContext(ctx, query, sku).Scan(
		&product.ID, &product.CategoryID, &product.SKU, &product.Name, &product.Description,
		&product.Price, &product.Stock, &product.ImageURL, &product.IsActive,
		&product.CreatedAt, &product.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return product, err
}

func (r *productRepository) Update(ctx context.Context, product *models.Product) error {
	query := `
		UPDATE products SET category_id = ?, sku = ?, name = ?, description = ?, price = ?, 
		       stock = ?, image_url = ?, is_active = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		product.CategoryID, product.SKU, product.Name, product.Description,
		product.Price, product.Stock, product.ImageURL, product.IsActive,
		product.UpdatedAt, product.ID,
	)
	return err
}

func (r *productRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM products WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *productRepository) UpdateStock(ctx context.Context, id string, quantity int) error {
	query := `UPDATE products SET stock = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, quantity, id)
	return err
}

func (r *productRepository) List(ctx context.Context, filter dto.ProductListFilter, pagination utils.Pagination) ([]*models.Product, int, error) {
	// Build where clause
	var conditions []string
	var args []interface{}

	if filter.CategoryID != "" {
		conditions = append(conditions, "p.category_id = ?")
		args = append(args, filter.CategoryID)
	}
	if filter.Search != "" {
		conditions = append(conditions, "(p.name LIKE ? OR p.sku LIKE ?)")
		searchTerm := "%" + filter.Search + "%"
		args = append(args, searchTerm, searchTerm)
	}
	if filter.MinPrice > 0 {
		conditions = append(conditions, "p.price >= ?")
		args = append(args, filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		conditions = append(conditions, "p.price <= ?")
		args = append(args, filter.MaxPrice)
	}
	if filter.InStock != nil && *filter.InStock {
		conditions = append(conditions, "p.stock > 0")
	}
	if filter.IsActive != nil {
		conditions = append(conditions, "p.is_active = ?")
		args = append(args, *filter.IsActive)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Get total count
	var total int
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM products p %s`, whereClause)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := fmt.Sprintf(`
		SELECT p.id, p.category_id, p.sku, p.name, COALESCE(p.description, ''), p.price, p.stock, COALESCE(p.image_url, ''), p.is_active, p.created_at, p.updated_at,
		       c.id, c.name, COALESCE(c.description, ''), c.slug, c.is_active, c.created_at, c.updated_at
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		%s
		ORDER BY p.%s
		LIMIT ? OFFSET ?
	`, whereClause, pagination.OrderBy())

	args = append(args, pagination.Limit(), pagination.Offset())
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		product := &models.Product{}
		var catID, catName, catDesc, catSlug sql.NullString
		var catIsActive sql.NullBool
		var catCreatedAt, catUpdatedAt sql.NullTime

		if err := rows.Scan(
			&product.ID, &product.CategoryID, &product.SKU, &product.Name, &product.Description,
			&product.Price, &product.Stock, &product.ImageURL, &product.IsActive,
			&product.CreatedAt, &product.UpdatedAt,
			&catID, &catName, &catDesc, &catSlug,
			&catIsActive, &catCreatedAt, &catUpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		if catID.Valid {
			product.Category = &models.Category{
				ID:          catID.String,
				Name:        catName.String,
				Description: catDesc.String,
				Slug:        catSlug.String,
				IsActive:    catIsActive.Bool,
				CreatedAt:   catCreatedAt.Time,
				UpdatedAt:   catUpdatedAt.Time,
			}
		}
		products = append(products, product)
	}

	return products, total, rows.Err()
}
