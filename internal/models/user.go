package models

import (
	"time"
)

// User represents a system user (admin, manager, or cashier)
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Name         string    `json:"name"`
	Role         string    `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserRole constants
const (
	RoleAdmin   = "admin"
	RoleManager = "manager"
	RoleCashier = "cashier"
)

// IsAdmin checks if the user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsManager checks if the user has manager role
func (u *User) IsManager() bool {
	return u.Role == RoleManager
}

// IsCashier checks if the user has cashier role
func (u *User) IsCashier() bool {
	return u.Role == RoleCashier
}

// CanManageProducts checks if the user can create/update/delete products
func (u *User) CanManageProducts() bool {
	return u.Role == RoleAdmin || u.Role == RoleManager
}

// CanViewReports checks if the user can view reports
func (u *User) CanViewReports() bool {
	return u.Role == RoleAdmin || u.Role == RoleManager
}

// CanDeleteRecords checks if the user can delete records
func (u *User) CanDeleteRecords() bool {
	return u.Role == RoleAdmin
}
