package dto

import "time"

// CreateProductRequest represents a request to create a product
type CreateProductRequest struct {
	CategoryID  string  `json:"category_id" validate:"required,uuid"`
	SKU         string  `json:"sku" validate:"required,min=3,max=50"`
	Name        string  `json:"name" validate:"required,min=2,max=200"`
	Description string  `json:"description" validate:"max=1000"`
	Price       float64 `json:"price" validate:"required,gte=0"`
	Stock       int     `json:"stock" validate:"gte=0"`
	ImageURL    string  `json:"image_url" validate:"omitempty,url"`
}

// UpdateProductRequest represents a request to update a product
type UpdateProductRequest struct {
	CategoryID  string  `json:"category_id" validate:"required,uuid"`
	SKU         string  `json:"sku" validate:"required,min=3,max=50"`
	Name        string  `json:"name" validate:"required,min=2,max=200"`
	Description string  `json:"description" validate:"max=1000"`
	Price       float64 `json:"price" validate:"required,gte=0"`
	Stock       int     `json:"stock" validate:"gte=0"`
	ImageURL    string  `json:"image_url" validate:"omitempty,url"`
	IsActive    *bool   `json:"is_active"`
}

// UpdateStockRequest represents a request to update product stock
type UpdateStockRequest struct {
	Quantity  int    `json:"quantity" validate:"required"`
	Operation string `json:"operation" validate:"required,oneof=add subtract set"`
}

// ProductResponse represents a product in responses
type ProductResponse struct {
	ID          string            `json:"id"`
	CategoryID  string            `json:"category_id"`
	SKU         string            `json:"sku"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Price       float64           `json:"price"`
	Stock       int               `json:"stock"`
	ImageURL    string            `json:"image_url,omitempty"`
	IsActive    bool              `json:"is_active"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Category    *CategoryResponse `json:"category,omitempty"`
}

// ProductListFilter represents filters for product listing
type ProductListFilter struct {
	CategoryID string  `form:"category_id"`
	Search     string  `form:"search"`
	MinPrice   float64 `form:"min_price"`
	MaxPrice   float64 `form:"max_price"`
	InStock    *bool   `form:"in_stock"`
	IsActive   *bool   `form:"is_active"`
}
