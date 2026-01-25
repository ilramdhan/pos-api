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
	"golang.org/x/crypto/bcrypt"
)

// UserService handles user management operations
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// List returns paginated list of users
func (s *UserService) List(ctx context.Context, role string, pagination utils.Pagination) ([]*dto.UserListResponse, int, error) {
	users, total, err := s.userRepo.List(ctx, role, pagination)
	if err != nil {
		return nil, 0, err
	}

	var result []*dto.UserListResponse
	for _, u := range users {
		result = append(result, &dto.UserListResponse{
			ID:          u.ID,
			Name:        u.Name,
			Email:       u.Email,
			Phone:       u.Phone,
			Role:        u.Role,
			IsActive:    u.IsActive,
			LastLoginAt: u.UpdatedAt.Format(time.RFC3339),
			CreatedAt:   u.CreatedAt.Format(time.RFC3339),
		})
	}

	return result, total, nil
}

// GetByID returns a user by ID
func (s *UserService) GetByID(ctx context.Context, id string) (*dto.UserListResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return &dto.UserListResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Phone:       user.Phone,
		Role:        user.Role,
		IsActive:    user.IsActive,
		LastLoginAt: user.UpdatedAt.Format(time.RFC3339),
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
	}, nil
}

// Create creates a new user
func (s *UserService) Create(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserListResponse, error) {
	// Check if email exists
	existing, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &models.User{
		ID:           uuid.New().String(),
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &dto.UserListResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}

// Update updates a user
func (s *UserService) Update(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserListResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return &dto.UserListResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}

// Delete deletes a user
func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}

// ResetPassword resets a user's password
func (s *UserService) ResetPassword(ctx context.Context, id, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(ctx, user)
}
