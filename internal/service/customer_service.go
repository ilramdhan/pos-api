package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ilramdhan/pos-api/internal/dto"
	"github.com/ilramdhan/pos-api/internal/models"
	"github.com/ilramdhan/pos-api/internal/repository"
	"github.com/ilramdhan/pos-api/internal/utils"
)

// CustomerService handles customer operations
type CustomerService struct {
	customerRepo repository.CustomerRepository
}

// NewCustomerService creates a new customer service
func NewCustomerService(customerRepo repository.CustomerRepository) *CustomerService {
	return &CustomerService{
		customerRepo: customerRepo,
	}
}

// Create creates a new customer
func (s *CustomerService) Create(ctx context.Context, req *dto.CreateCustomerRequest) (*dto.CustomerResponse, error) {
	now := time.Now()
	customer := &models.Customer{
		ID:            uuid.New().String(),
		Name:          req.Name,
		Email:         req.Email,
		Phone:         req.Phone,
		Address:       req.Address,
		LoyaltyPoints: 0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.customerRepo.Create(ctx, customer); err != nil {
		return nil, err
	}

	return s.toResponse(customer), nil
}

// GetByID retrieves a customer by ID
func (s *CustomerService) GetByID(ctx context.Context, id string) (*dto.CustomerResponse, error) {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, errors.New("customer not found")
	}

	return s.toResponse(customer), nil
}

// Update updates a customer
func (s *CustomerService) Update(ctx context.Context, id string, req *dto.UpdateCustomerRequest) (*dto.CustomerResponse, error) {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, errors.New("customer not found")
	}

	customer.Name = req.Name
	customer.Email = req.Email
	customer.Phone = req.Phone
	customer.Address = req.Address
	customer.UpdatedAt = time.Now()

	if err := s.customerRepo.Update(ctx, customer); err != nil {
		return nil, err
	}

	return s.toResponse(customer), nil
}

// Delete deletes a customer
func (s *CustomerService) Delete(ctx context.Context, id string) error {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if customer == nil {
		return errors.New("customer not found")
	}

	return s.customerRepo.Delete(ctx, id)
}

// UpdateLoyaltyPoints updates a customer's loyalty points
func (s *CustomerService) UpdateLoyaltyPoints(ctx context.Context, id string, req *dto.UpdateLoyaltyPointsRequest) (*dto.CustomerResponse, error) {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, errors.New("customer not found")
	}

	var newPoints int
	switch req.Operation {
	case "add":
		newPoints = customer.LoyaltyPoints + req.Points
	case "deduct":
		newPoints = customer.LoyaltyPoints - req.Points
		if newPoints < 0 {
			return nil, errors.New("insufficient loyalty points")
		}
	default:
		return nil, errors.New("invalid operation")
	}

	if err := s.customerRepo.UpdateLoyaltyPoints(ctx, id, newPoints); err != nil {
		return nil, err
	}

	customer.LoyaltyPoints = newPoints
	return s.toResponse(customer), nil
}

// List lists customers with pagination and filters
func (s *CustomerService) List(ctx context.Context, filter dto.CustomerListFilter, pagination utils.Pagination) ([]*dto.CustomerResponse, int, error) {
	customers, total, err := s.customerRepo.List(ctx, filter, pagination)
	if err != nil {
		return nil, 0, err
	}

	var responses []*dto.CustomerResponse
	for _, customer := range customers {
		responses = append(responses, s.toResponse(customer))
	}

	return responses, total, nil
}

func (s *CustomerService) toResponse(customer *models.Customer) *dto.CustomerResponse {
	return &dto.CustomerResponse{
		ID:            customer.ID,
		Name:          customer.Name,
		Email:         customer.Email,
		Phone:         customer.Phone,
		Address:       customer.Address,
		LoyaltyPoints: customer.LoyaltyPoints,
		CreatedAt:     customer.CreatedAt,
		UpdatedAt:     customer.UpdatedAt,
	}
}
