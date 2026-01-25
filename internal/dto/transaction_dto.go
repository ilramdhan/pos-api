package dto

import "time"

// CreateTransactionRequest represents a request to create a transaction
type CreateTransactionRequest struct {
	CustomerID     *string                    `json:"customer_id" validate:"omitempty,uuid"`
	PaymentMethod  string                     `json:"payment_method" validate:"required,oneof=cash card ewallet"`
	DiscountAmount float64                    `json:"discount_amount" validate:"gte=0"`
	Notes          string                     `json:"notes" validate:"max=500"`
	Items          []CreateTransactionItemDTO `json:"items" validate:"required,min=1,dive"`
}

// CreateTransactionItemDTO represents a line item in a transaction request
type CreateTransactionItemDTO struct {
	ProductID string `json:"product_id" validate:"required,uuid"`
	Quantity  int    `json:"quantity" validate:"required,gt=0"`
}

// UpdateTransactionStatusRequest represents a request to update transaction status
type UpdateTransactionStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=completed cancelled refunded"`
}

// TransactionResponse represents a transaction in responses
type TransactionResponse struct {
	ID             string                    `json:"id"`
	UserID         string                    `json:"user_id"`
	CustomerID     *string                   `json:"customer_id,omitempty"`
	InvoiceNumber  string                    `json:"invoice_number"`
	Subtotal       float64                   `json:"subtotal"`
	TaxAmount      float64                   `json:"tax_amount"`
	DiscountAmount float64                   `json:"discount_amount"`
	TotalAmount    float64                   `json:"total_amount"`
	PaymentMethod  string                    `json:"payment_method"`
	Status         string                    `json:"status"`
	Notes          string                    `json:"notes,omitempty"`
	CreatedAt      time.Time                 `json:"created_at"`
	UpdatedAt      time.Time                 `json:"updated_at"`
	User           *UserResponse             `json:"user,omitempty"`
	Customer       *CustomerResponse         `json:"customer,omitempty"`
	Items          []TransactionItemResponse `json:"items,omitempty"`
}

// TransactionItemResponse represents a transaction item in responses
type TransactionItemResponse struct {
	ID          string  `json:"id"`
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	UnitPrice   float64 `json:"unit_price"`
	Quantity    int     `json:"quantity"`
	Subtotal    float64 `json:"subtotal"`
}

// TransactionListFilter represents filters for transaction listing
type TransactionListFilter struct {
	UserID        string `form:"user_id"`
	CustomerID    string `form:"customer_id"`
	Status        string `form:"status"`
	PaymentMethod string `form:"payment_method"`
	DateFrom      string `form:"date_from"`
	DateTo        string `form:"date_to"`
}

// SalesReportRequest represents a request for sales report
type SalesReportRequest struct {
	DateFrom string `form:"date_from" validate:"required"`
	DateTo   string `form:"date_to" validate:"required"`
}

// DailySalesReport represents daily sales summary
type DailySalesReport struct {
	Date              string  `json:"date"`
	TotalTransactions int     `json:"total_transactions"`
	TotalAmount       float64 `json:"total_amount"`
	TotalItems        int     `json:"total_items"`
}

// MonthlySalesReport represents monthly sales summary
type MonthlySalesReport struct {
	Month             string  `json:"month"`
	TotalTransactions int     `json:"total_transactions"`
	TotalAmount       float64 `json:"total_amount"`
	TotalItems        int     `json:"total_items"`
}

// TopProductReport represents top selling products
type TopProductReport struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	SKU         string  `json:"sku"`
	TotalSold   int     `json:"total_sold"`
	TotalAmount float64 `json:"total_amount"`
}
