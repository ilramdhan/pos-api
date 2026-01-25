package models

import (
	"time"
)

// Transaction represents a sales transaction
type Transaction struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	CustomerID     *string   `json:"customer_id,omitempty"`
	InvoiceNumber  string    `json:"invoice_number"`
	Subtotal       float64   `json:"subtotal"`
	TaxAmount      float64   `json:"tax_amount"`
	DiscountAmount float64   `json:"discount_amount"`
	TotalAmount    float64   `json:"total_amount"`
	PaymentMethod  string    `json:"payment_method"`
	Status         string    `json:"status"`
	Notes          string    `json:"notes,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Joined fields
	User     *User             `json:"user,omitempty"`
	Customer *Customer         `json:"customer,omitempty"`
	Items    []TransactionItem `json:"items,omitempty"`
}

// TransactionItem represents a line item in a transaction
type TransactionItem struct {
	ID            string    `json:"id"`
	TransactionID string    `json:"transaction_id"`
	ProductID     string    `json:"product_id"`
	ProductName   string    `json:"product_name"`
	UnitPrice     float64   `json:"unit_price"`
	Quantity      int       `json:"quantity"`
	Subtotal      float64   `json:"subtotal"`
	CreatedAt     time.Time `json:"created_at"`

	// Joined fields
	Product *Product `json:"product,omitempty"`
}

// Transaction status constants
const (
	StatusPending   = "pending"
	StatusCompleted = "completed"
	StatusCancelled = "cancelled"
	StatusRefunded  = "refunded"
)

// Payment method constants
const (
	PaymentCash    = "cash"
	PaymentCard    = "card"
	PaymentEWallet = "ewallet"
)

// IsCompleted checks if the transaction is completed
func (t *Transaction) IsCompleted() bool {
	return t.Status == StatusCompleted
}

// IsCancellable checks if the transaction can be cancelled
func (t *Transaction) IsCancellable() bool {
	return t.Status == StatusPending
}

// IsRefundable checks if the transaction can be refunded
func (t *Transaction) IsRefundable() bool {
	return t.Status == StatusCompleted
}

// CalculateTotals calculates and sets the transaction totals
func (t *Transaction) CalculateTotals(taxRate float64) {
	var subtotal float64
	for _, item := range t.Items {
		subtotal += item.Subtotal
	}
	t.Subtotal = subtotal
	t.TaxAmount = subtotal * taxRate
	t.TotalAmount = subtotal + t.TaxAmount - t.DiscountAmount
}
