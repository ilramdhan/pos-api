package dto

// CreateUserRequest represents a create user request
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Phone    string `json:"phone" validate:"omitempty"`
	Role     string `json:"role" validate:"required,oneof=admin manager cashier"`
}

// UpdateUserRequest represents an update user request
type UpdateUserRequest struct {
	Name     string `json:"name" validate:"omitempty,min=2,max=100"`
	Phone    string `json:"phone" validate:"omitempty"`
	Role     string `json:"role" validate:"omitempty,oneof=admin manager cashier"`
	IsActive *bool  `json:"is_active" validate:"omitempty"`
}

// ResetUserPasswordRequest represents a password reset request
type ResetUserPasswordRequest struct {
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// UserListResponse represents a user in list response
type UserListResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone,omitempty"`
	Role        string `json:"role"`
	IsActive    bool   `json:"is_active"`
	LastLoginAt string `json:"last_login_at,omitempty"`
	CreatedAt   string `json:"created_at"`
}
