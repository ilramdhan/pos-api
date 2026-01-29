package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ilramdhan/pos-api/internal/models"
	"github.com/ilramdhan/pos-api/internal/utils"
)

type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, name, phone, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.Name, user.Phone, user.Role,
		user.IsActive, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, name, COALESCE(phone, '') as phone, role, is_active, created_at, updated_at
		FROM users WHERE id = $1
	`
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Phone, &user.Role,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, name, COALESCE(phone, '') as phone, role, is_active, created_at, updated_at
		FROM users WHERE email = $1
	`
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Phone, &user.Role,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users SET email = $1, name = $2, phone = $3, role = $4, is_active = $5, updated_at = $6
		WHERE id = $7
	`
	_, err := r.db.ExecContext(ctx, query,
		user.Email, user.Name, user.Phone, user.Role, user.IsActive, user.UpdatedAt, user.ID,
	)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *userRepository) List(ctx context.Context, role string, pagination utils.Pagination) ([]*models.User, int, error) {
	// Build where clause
	whereClause := ""
	var args []interface{}
	argIndex := 1

	if role != "" {
		whereClause = fmt.Sprintf("WHERE role = $%d", argIndex)
		args = append(args, role)
		argIndex++
	}

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM users ` + whereClause
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := fmt.Sprintf(`
		SELECT id, email, password_hash, name, COALESCE(phone, '') as phone, role, is_active, created_at, updated_at
		FROM users %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereClause, pagination.OrderBy(), argIndex, argIndex+1)

	args = append(args, pagination.Limit(), pagination.Offset())
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Phone, &user.Role,
			&user.IsActive, &user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, rows.Err()
}
