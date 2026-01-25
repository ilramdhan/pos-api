package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilramdhan/pos-api/internal/dto"
	"github.com/ilramdhan/pos-api/internal/middleware"
	"github.com/ilramdhan/pos-api/internal/service"
	"github.com/ilramdhan/pos-api/internal/utils"
)

// UserHandler handles user management endpoints
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// List handles GET /api/v1/users
func (h *UserHandler) List(c *gin.Context) {
	pagination := utils.GetPagination(c)
	role := c.Query("role")

	users, total, err := h.userService.List(c.Request.Context(), role, pagination)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithMeta(c, "Users retrieved successfully", users, &utils.Meta{
		Page:       pagination.Page,
		PerPage:    pagination.PerPage,
		Total:      total,
		TotalPages: (total + pagination.PerPage - 1) / pagination.PerPage,
	})
}

// Get handles GET /api/v1/users/:id
func (h *UserHandler) Get(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User retrieved successfully", user)
}

// Create handles POST /api/v1/users
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if errors, ok := utils.Validate(&req); !ok {
		utils.ValidationErrorResponse(c, errors)
		return
	}

	user, err := h.userService.Create(c.Request.Context(), &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.CreatedResponse(c, "User created successfully", user)
}

// Update handles PUT /api/v1/users/:id
func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	user, err := h.userService.Update(c.Request.Context(), id, &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User updated successfully", user)
}

// Delete handles DELETE /api/v1/users/:id
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	// Prevent deleting self
	claims := middleware.GetCurrentUser(c)
	if claims != nil && claims.UserID == id {
		utils.BadRequest(c, "Cannot delete your own account")
		return
	}

	if err := h.userService.Delete(c.Request.Context(), id); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User deleted successfully", nil)
}

// ResetPassword handles PUT /api/v1/users/:id/reset-password
func (h *UserHandler) ResetPassword(c *gin.Context) {
	id := c.Param("id")

	var req dto.ResetUserPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.userService.ResetPassword(c.Request.Context(), id, req.NewPassword); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Password reset successfully", nil)
}
