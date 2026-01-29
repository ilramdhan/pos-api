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
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
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
		WHERE p.id = $1
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
		FROM products WHERE sku = $1
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
		UPDATE products SET category_id = $1, sku = $2, name = $3, description = $4, price = $5, 
		       stock = $6, image_url = $7, is_active = $8, updated_at = $9
		WHERE id = $10
	`
	_, err := r.db.ExecContext(ctx, query,
		product.CategoryID, product.SKU, product.Name, product.Description,
		product.Price, product.Stock, product.ImageURL, product.IsActive,
		product.UpdatedAt, product.ID,
	)
	return err
}

func (r *productRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *productRepository) UpdateStock(ctx context.Context, id string, quantity int) error {
	query := `UPDATE products SET stock = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, quantity, id)
	return err
}

func (r *productRepository) List(ctx context.Context, filter dto.ProductListFilter, pagination utils.Pagination) ([]*models.Product, int, error) {
	// Build where clause
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.CategoryID != "" {
		conditions = append(conditions, fmt.Sprintf("p.category_id = $%d", argIndex))
		args = append(args, filter.CategoryID)
		argIndex++
	}
	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(p.name ILIKE $%d OR p.sku ILIKE $%d)", argIndex, argIndex+1))
		searchTerm := "%" + filter.Search + "%"
		args = append(args, searchTerm, searchTerm)
		argIndex += 2
	}
	if filter.MinPrice > 0 {
		conditions = append(conditions, fmt.Sprintf("p.price >= $%d", argIndex))
		args = append(args, filter.MinPrice)
		argIndex++
	}
	if filter.MaxPrice > 0 {
		conditions = append(conditions, fmt.Sprintf("p.price <= $%d", argIndex))
		args = append(args, filter.MaxPrice)
		argIndex++
	}
	if filter.InStock != nil && *filter.InStock {
		conditions = append(conditions, "p.stock > 0")
	}
	if filter.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("p.is_active = $%d", argIndex))
		args = append(args, *filter.IsActive)
		argIndex++
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

	// Build paginated query
	query := fmt.Sprintf(`
		SELECT p.id, p.category_id, p.sku, p.name, COALESCE(p.description, ''), p.price, p.stock, COALESCE(p.image_url, ''), p.is_active, p.created_at, p.updated_at,
		       c.id, c.name, COALESCE(c.description, ''), c.slug, c.is_active, c.created_at, c.updated_at
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		%s
		ORDER BY p.%s
		LIMIT $%d OFFSET $%d
	`, whereClause, pagination.OrderBy(), argIndex, argIndex+1)

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
