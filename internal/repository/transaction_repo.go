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

type transactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, transaction *models.Transaction) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert transaction
	query := `
		INSERT INTO transactions (id, user_id, customer_id, invoice_number, subtotal, tax_amount, 
		                         discount_amount, total_amount, payment_method, status, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = tx.ExecContext(ctx, query,
		transaction.ID, transaction.UserID, transaction.CustomerID, transaction.InvoiceNumber,
		transaction.Subtotal, transaction.TaxAmount, transaction.DiscountAmount, transaction.TotalAmount,
		transaction.PaymentMethod, transaction.Status, transaction.Notes,
		transaction.CreatedAt, transaction.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Insert transaction items
	itemQuery := `
		INSERT INTO transaction_items (id, transaction_id, product_id, product_name, unit_price, quantity, subtotal, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	for _, item := range transaction.Items {
		_, err = tx.ExecContext(ctx, itemQuery,
			item.ID, item.TransactionID, item.ProductID, item.ProductName,
			item.UnitPrice, item.Quantity, item.Subtotal, item.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *transactionRepository) GetByID(ctx context.Context, id string) (*models.Transaction, error) {
	query := `
		SELECT t.id, t.user_id, t.customer_id, t.invoice_number, t.subtotal, t.tax_amount,
		       t.discount_amount, t.total_amount, t.payment_method, t.status, t.notes, t.created_at, t.updated_at,
		       u.id, u.email, u.name, u.role, u.is_active
		FROM transactions t
		LEFT JOIN users u ON t.user_id = u.id
		WHERE t.id = ?
	`

	transaction := &models.Transaction{}
	user := &models.User{}
	var customerID sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&transaction.ID, &transaction.UserID, &customerID, &transaction.InvoiceNumber,
		&transaction.Subtotal, &transaction.TaxAmount, &transaction.DiscountAmount,
		&transaction.TotalAmount, &transaction.PaymentMethod, &transaction.Status,
		&transaction.Notes, &transaction.CreatedAt, &transaction.UpdatedAt,
		&user.ID, &user.Email, &user.Name, &user.Role, &user.IsActive,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if customerID.Valid {
		transaction.CustomerID = &customerID.String
	}
	transaction.User = user

	// Get transaction items
	itemQuery := `
		SELECT id, transaction_id, product_id, product_name, unit_price, quantity, subtotal, created_at
		FROM transaction_items WHERE transaction_id = ?
	`
	rows, err := r.db.QueryContext(ctx, itemQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := models.TransactionItem{}
		if err := rows.Scan(
			&item.ID, &item.TransactionID, &item.ProductID, &item.ProductName,
			&item.UnitPrice, &item.Quantity, &item.Subtotal, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		transaction.Items = append(transaction.Items, item)
	}

	return transaction, rows.Err()
}

func (r *transactionRepository) GetByInvoiceNumber(ctx context.Context, invoiceNumber string) (*models.Transaction, error) {
	query := `
		SELECT id, user_id, customer_id, invoice_number, subtotal, tax_amount,
		       discount_amount, total_amount, payment_method, status, notes, created_at, updated_at
		FROM transactions WHERE invoice_number = ?
	`
	transaction := &models.Transaction{}
	var customerID sql.NullString

	err := r.db.QueryRowContext(ctx, query, invoiceNumber).Scan(
		&transaction.ID, &transaction.UserID, &customerID, &transaction.InvoiceNumber,
		&transaction.Subtotal, &transaction.TaxAmount, &transaction.DiscountAmount,
		&transaction.TotalAmount, &transaction.PaymentMethod, &transaction.Status,
		&transaction.Notes, &transaction.CreatedAt, &transaction.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if customerID.Valid {
		transaction.CustomerID = &customerID.String
	}
	return transaction, err
}

func (r *transactionRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE transactions SET status = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *transactionRepository) List(ctx context.Context, filter dto.TransactionListFilter, pagination utils.Pagination) ([]*models.Transaction, int, error) {
	var conditions []string
	var args []interface{}

	if filter.UserID != "" {
		conditions = append(conditions, "t.user_id = ?")
		args = append(args, filter.UserID)
	}
	if filter.CustomerID != "" {
		conditions = append(conditions, "t.customer_id = ?")
		args = append(args, filter.CustomerID)
	}
	if filter.Status != "" {
		conditions = append(conditions, "t.status = ?")
		args = append(args, filter.Status)
	}
	if filter.PaymentMethod != "" {
		conditions = append(conditions, "t.payment_method = ?")
		args = append(args, filter.PaymentMethod)
	}
	if filter.DateFrom != "" {
		conditions = append(conditions, "DATE(t.created_at) >= ?")
		args = append(args, filter.DateFrom)
	}
	if filter.DateTo != "" {
		conditions = append(conditions, "DATE(t.created_at) <= ?")
		args = append(args, filter.DateTo)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Get total count
	var total int
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM transactions t %s`, whereClause)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := fmt.Sprintf(`
		SELECT t.id, t.user_id, t.customer_id, t.invoice_number, t.subtotal, t.tax_amount,
		       t.discount_amount, t.total_amount, t.payment_method, t.status, t.notes, t.created_at, t.updated_at
		FROM transactions t
		%s
		ORDER BY t.%s
		LIMIT ? OFFSET ?
	`, whereClause, pagination.OrderBy())

	args = append(args, pagination.Limit(), pagination.Offset())
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		transaction := &models.Transaction{}
		var customerID sql.NullString
		if err := rows.Scan(
			&transaction.ID, &transaction.UserID, &customerID, &transaction.InvoiceNumber,
			&transaction.Subtotal, &transaction.TaxAmount, &transaction.DiscountAmount,
			&transaction.TotalAmount, &transaction.PaymentMethod, &transaction.Status,
			&transaction.Notes, &transaction.CreatedAt, &transaction.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		if customerID.Valid {
			transaction.CustomerID = &customerID.String
		}
		transactions = append(transactions, transaction)
	}

	return transactions, total, rows.Err()
}

func (r *transactionRepository) GetDailySales(ctx context.Context, dateFrom, dateTo string) ([]dto.DailySalesReport, error) {
	query := `
		SELECT DATE(t.created_at) as date, 
		       COUNT(*) as total_transactions,
		       SUM(t.total_amount) as total_amount,
		       SUM(ti.quantity) as total_items
		FROM transactions t
		LEFT JOIN transaction_items ti ON t.id = ti.transaction_id
		WHERE t.status = 'completed'
		  AND DATE(t.created_at) >= ? AND DATE(t.created_at) <= ?
		GROUP BY DATE(t.created_at)
		ORDER BY date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []dto.DailySalesReport
	for rows.Next() {
		var report dto.DailySalesReport
		var totalItems sql.NullInt64
		if err := rows.Scan(&report.Date, &report.TotalTransactions, &report.TotalAmount, &totalItems); err != nil {
			return nil, err
		}
		if totalItems.Valid {
			report.TotalItems = int(totalItems.Int64)
		}
		reports = append(reports, report)
	}

	return reports, rows.Err()
}

func (r *transactionRepository) GetMonthlySales(ctx context.Context, dateFrom, dateTo string) ([]dto.MonthlySalesReport, error) {
	query := `
		SELECT strftime('%Y-%m', t.created_at) as month,
		       COUNT(*) as total_transactions,
		       SUM(t.total_amount) as total_amount,
		       SUM(ti.quantity) as total_items
		FROM transactions t
		LEFT JOIN transaction_items ti ON t.id = ti.transaction_id
		WHERE t.status = 'completed'
		  AND DATE(t.created_at) >= ? AND DATE(t.created_at) <= ?
		GROUP BY strftime('%Y-%m', t.created_at)
		ORDER BY month DESC
	`

	rows, err := r.db.QueryContext(ctx, query, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []dto.MonthlySalesReport
	for rows.Next() {
		var report dto.MonthlySalesReport
		var totalItems sql.NullInt64
		if err := rows.Scan(&report.Month, &report.TotalTransactions, &report.TotalAmount, &totalItems); err != nil {
			return nil, err
		}
		if totalItems.Valid {
			report.TotalItems = int(totalItems.Int64)
		}
		reports = append(reports, report)
	}

	return reports, rows.Err()
}

func (r *transactionRepository) GetTopProducts(ctx context.Context, limit int, dateFrom, dateTo string) ([]dto.TopProductReport, error) {
	query := `
		SELECT ti.product_id, ti.product_name, p.sku,
		       SUM(ti.quantity) as total_sold,
		       SUM(ti.subtotal) as total_amount
		FROM transaction_items ti
		JOIN transactions t ON ti.transaction_id = t.id
		LEFT JOIN products p ON ti.product_id = p.id
		WHERE t.status = 'completed'
		  AND DATE(t.created_at) >= ? AND DATE(t.created_at) <= ?
		GROUP BY ti.product_id, ti.product_name, p.sku
		ORDER BY total_sold DESC
		LIMIT ?
	`

	rows, err := r.db.QueryContext(ctx, query, dateFrom, dateTo, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []dto.TopProductReport
	for rows.Next() {
		var report dto.TopProductReport
		var sku sql.NullString
		if err := rows.Scan(&report.ProductID, &report.ProductName, &sku, &report.TotalSold, &report.TotalAmount); err != nil {
			return nil, err
		}
		if sku.Valid {
			report.SKU = sku.String
		}
		reports = append(reports, report)
	}

	return reports, rows.Err()
}
