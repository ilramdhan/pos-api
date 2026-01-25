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

type customerRepository struct {
	db *sql.DB
}

// NewCustomerRepository creates a new customer repository
func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(ctx context.Context, customer *models.Customer) error {
	query := `
		INSERT INTO customers (id, name, email, phone, address, loyalty_points, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		customer.ID, customer.Name, customer.Email, customer.Phone, customer.Address,
		customer.LoyaltyPoints, customer.CreatedAt, customer.UpdatedAt,
	)
	return err
}

func (r *customerRepository) GetByID(ctx context.Context, id string) (*models.Customer, error) {
	query := `
		SELECT id, name, email, phone, address, loyalty_points, created_at, updated_at
		FROM customers WHERE id = ?
	`
	customer := &models.Customer{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&customer.ID, &customer.Name, &customer.Email, &customer.Phone, &customer.Address,
		&customer.LoyaltyPoints, &customer.CreatedAt, &customer.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return customer, err
}

func (r *customerRepository) Update(ctx context.Context, customer *models.Customer) error {
	query := `
		UPDATE customers SET name = ?, email = ?, phone = ?, address = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		customer.Name, customer.Email, customer.Phone, customer.Address,
		customer.UpdatedAt, customer.ID,
	)
	return err
}

func (r *customerRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM customers WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *customerRepository) UpdateLoyaltyPoints(ctx context.Context, id string, points int) error {
	query := `UPDATE customers SET loyalty_points = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, points, id)
	return err
}

func (r *customerRepository) List(ctx context.Context, filter dto.CustomerListFilter, pagination utils.Pagination) ([]*models.Customer, int, error) {
	// Build where clause
	var conditions []string
	var args []interface{}

	if filter.Search != "" {
		conditions = append(conditions, "(name LIKE ? OR email LIKE ? OR phone LIKE ?)")
		searchTerm := "%" + filter.Search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Get total count
	var total int
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM customers %s`, whereClause)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := fmt.Sprintf(`
		SELECT id, name, email, phone, address, loyalty_points, created_at, updated_at
		FROM customers
		%s
		ORDER BY %s
		LIMIT ? OFFSET ?
	`, whereClause, pagination.OrderBy())

	args = append(args, pagination.Limit(), pagination.Offset())
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var customers []*models.Customer
	for rows.Next() {
		customer := &models.Customer{}
		if err := rows.Scan(
			&customer.ID, &customer.Name, &customer.Email, &customer.Phone, &customer.Address,
			&customer.LoyaltyPoints, &customer.CreatedAt, &customer.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		customers = append(customers, customer)
	}

	return customers, total, rows.Err()
}
