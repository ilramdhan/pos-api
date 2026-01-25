package models

import (
	"time"
)

// Product represents a product in the inventory
type Product struct {
	ID          string    `json:"id"`
	CategoryID  string    `json:"category_id"`
	SKU         string    `json:"sku"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	ImageURL    string    `json:"image_url,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Joined fields
	Category *Category `json:"category,omitempty"`
}

// IsInStock checks if the product has available stock
func (p *Product) IsInStock() bool {
	return p.Stock > 0
}

// HasSufficientStock checks if there's enough stock for a given quantity
func (p *Product) HasSufficientStock(quantity int) bool {
	return p.Stock >= quantity
}
