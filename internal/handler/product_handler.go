package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilramdhan/pos-api/internal/dto"
	"github.com/ilramdhan/pos-api/internal/service"
	"github.com/ilramdhan/pos-api/internal/utils"
)

// ProductHandler handles product endpoints
type ProductHandler struct {
	productService *service.ProductService
}

// NewProductHandler creates a new product handler
func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// List handles GET /api/v1/products
func (h *ProductHandler) List(c *gin.Context) {
	pagination := utils.GetPagination(c)

	var filter dto.ProductListFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.BadRequest(c, "Invalid query parameters")
		return
	}

	products, total, err := h.productService.List(c.Request.Context(), filter, pagination)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	meta := utils.NewMeta(pagination.Page, pagination.PerPage, total)
	utils.SuccessWithMeta(c, "Products retrieved successfully", products, meta)
}

// Get handles GET /api/v1/products/:id
func (h *ProductHandler) Get(c *gin.Context) {
	id := c.Param("id")

	product, err := h.productService.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product retrieved successfully", product)
}

// Create handles POST /api/v1/products
func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if errors, ok := utils.Validate(&req); !ok {
		utils.ValidationErrorResponse(c, errors)
		return
	}

	product, err := h.productService.Create(c.Request.Context(), &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.CreatedResponse(c, "Product created successfully", product)
}

// Update handles PUT /api/v1/products/:id
func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if errors, ok := utils.Validate(&req); !ok {
		utils.ValidationErrorResponse(c, errors)
		return
	}

	product, err := h.productService.Update(c.Request.Context(), id, &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product updated successfully", product)
}

// Delete handles DELETE /api/v1/products/:id
func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.productService.Delete(c.Request.Context(), id); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product deleted successfully", nil)
}

// UpdateStock handles PATCH /api/v1/products/:id/stock
func (h *ProductHandler) UpdateStock(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if errors, ok := utils.Validate(&req); !ok {
		utils.ValidationErrorResponse(c, errors)
		return
	}

	product, err := h.productService.UpdateStock(c.Request.Context(), id, &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Stock updated successfully", product)
}
