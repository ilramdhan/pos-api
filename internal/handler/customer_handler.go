package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilramdhan/pos-api/internal/dto"
	"github.com/ilramdhan/pos-api/internal/service"
	"github.com/ilramdhan/pos-api/internal/utils"
)

// CustomerHandler handles customer endpoints
type CustomerHandler struct {
	customerService *service.CustomerService
}

// NewCustomerHandler creates a new customer handler
func NewCustomerHandler(customerService *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{customerService: customerService}
}

// List handles GET /api/v1/customers
func (h *CustomerHandler) List(c *gin.Context) {
	pagination := utils.GetPagination(c)

	var filter dto.CustomerListFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.BadRequest(c, "Invalid query parameters")
		return
	}

	customers, total, err := h.customerService.List(c.Request.Context(), filter, pagination)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	meta := utils.NewMeta(pagination.Page, pagination.PerPage, total)
	utils.SuccessWithMeta(c, "Customers retrieved successfully", customers, meta)
}

// Get handles GET /api/v1/customers/:id
func (h *CustomerHandler) Get(c *gin.Context) {
	id := c.Param("id")

	customer, err := h.customerService.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Customer retrieved successfully", customer)
}

// Create handles POST /api/v1/customers
func (h *CustomerHandler) Create(c *gin.Context) {
	var req dto.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if errors, ok := utils.Validate(&req); !ok {
		utils.ValidationErrorResponse(c, errors)
		return
	}

	customer, err := h.customerService.Create(c.Request.Context(), &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.CreatedResponse(c, "Customer created successfully", customer)
}

// Update handles PUT /api/v1/customers/:id
func (h *CustomerHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if errors, ok := utils.Validate(&req); !ok {
		utils.ValidationErrorResponse(c, errors)
		return
	}

	customer, err := h.customerService.Update(c.Request.Context(), id, &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Customer updated successfully", customer)
}

// Delete handles DELETE /api/v1/customers/:id
func (h *CustomerHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.customerService.Delete(c.Request.Context(), id); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Customer deleted successfully", nil)
}
