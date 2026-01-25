package models

import (
	"time"
)

// Customer represents a customer in the system
type Customer struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email,omitempty"`
	Phone         string    `json:"phone,omitempty"`
	Address       string    `json:"address,omitempty"`
	LoyaltyPoints int       `json:"loyalty_points"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// AddLoyaltyPoints adds points to the customer
func (c *Customer) AddLoyaltyPoints(points int) {
	c.LoyaltyPoints += points
}

// DeductLoyaltyPoints deducts points from the customer, returns false if insufficient
func (c *Customer) DeductLoyaltyPoints(points int) bool {
	if c.LoyaltyPoints < points {
		return false
	}
	c.LoyaltyPoints -= points
	return true
}
