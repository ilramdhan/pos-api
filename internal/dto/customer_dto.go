package dto

import "time"

// CreateCustomerRequest represents a request to create a customer
type CreateCustomerRequest struct {
	Name    string `json:"name" validate:"required,min=2,max=100"`
	Email   string `json:"email" validate:"omitempty,email"`
	Phone   string `json:"phone" validate:"omitempty,min=10,max=20"`
	Address string `json:"address" validate:"max=500"`
}

// UpdateCustomerRequest represents a request to update a customer
type UpdateCustomerRequest struct {
	Name    string `json:"name" validate:"required,min=2,max=100"`
	Email   string `json:"email" validate:"omitempty,email"`
	Phone   string `json:"phone" validate:"omitempty,min=10,max=20"`
	Address string `json:"address" validate:"max=500"`
}

// UpdateLoyaltyPointsRequest represents a request to update loyalty points
type UpdateLoyaltyPointsRequest struct {
	Points    int    `json:"points" validate:"required,gt=0"`
	Operation string `json:"operation" validate:"required,oneof=add deduct"`
}

// CustomerResponse represents a customer in responses
type CustomerResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email,omitempty"`
	Phone         string    `json:"phone,omitempty"`
	Address       string    `json:"address,omitempty"`
	LoyaltyPoints int       `json:"loyalty_points"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CustomerListFilter represents filters for customer listing
type CustomerListFilter struct {
	Search string `form:"search"`
}
