package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilramdhan/pos-api/internal/dto"
	"github.com/ilramdhan/pos-api/internal/service"
	"github.com/ilramdhan/pos-api/internal/utils"
)

// CategoryHandler handles category endpoints
type CategoryHandler struct {
	categoryService *service.CategoryService
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

// List handles GET /api/v1/categories
func (h *CategoryHandler) List(c *gin.Context) {
	pagination := utils.GetPagination(c)

	categories, total, err := h.categoryService.List(c.Request.Context(), pagination)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	meta := utils.NewMeta(pagination.Page, pagination.PerPage, total)
	utils.SuccessWithMeta(c, "Categories retrieved successfully", categories, meta)
}

// Get handles GET /api/v1/categories/:id
func (h *CategoryHandler) Get(c *gin.Context) {
	id := c.Param("id")

	category, err := h.categoryService.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Category retrieved successfully", category)
}

// Create handles POST /api/v1/categories
func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if errors, ok := utils.Validate(&req); !ok {
		utils.ValidationErrorResponse(c, errors)
		return
	}

	category, err := h.categoryService.Create(c.Request.Context(), &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.CreatedResponse(c, "Category created successfully", category)
}

// Update handles PUT /api/v1/categories/:id
func (h *CategoryHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if errors, ok := utils.Validate(&req); !ok {
		utils.ValidationErrorResponse(c, errors)
		return
	}

	category, err := h.categoryService.Update(c.Request.Context(), id, &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Category updated successfully", category)
}

// Delete handles DELETE /api/v1/categories/:id
func (h *CategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.categoryService.Delete(c.Request.Context(), id); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Category deleted successfully", nil)
}
