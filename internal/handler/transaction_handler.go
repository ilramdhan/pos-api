package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilramdhan/pos-api/internal/dto"
	"github.com/ilramdhan/pos-api/internal/middleware"
	"github.com/ilramdhan/pos-api/internal/service"
	"github.com/ilramdhan/pos-api/internal/utils"
)

// TransactionHandler handles transaction endpoints
type TransactionHandler struct {
	transactionService *service.TransactionService
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(transactionService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

// List handles GET /api/v1/transactions
func (h *TransactionHandler) List(c *gin.Context) {
	pagination := utils.GetPagination(c)

	var filter dto.TransactionListFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.BadRequest(c, "Invalid query parameters")
		return
	}

	transactions, total, err := h.transactionService.List(c.Request.Context(), filter, pagination)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	meta := utils.NewMeta(pagination.Page, pagination.PerPage, total)
	utils.SuccessWithMeta(c, "Transactions retrieved successfully", transactions, meta)
}

// Get handles GET /api/v1/transactions/:id
func (h *TransactionHandler) Get(c *gin.Context) {
	id := c.Param("id")

	transaction, err := h.transactionService.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Transaction retrieved successfully", transaction)
}

// Create handles POST /api/v1/transactions
func (h *TransactionHandler) Create(c *gin.Context) {
	claims := middleware.GetCurrentUser(c)
	if claims == nil {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if errors, ok := utils.Validate(&req); !ok {
		utils.ValidationErrorResponse(c, errors)
		return
	}

	transaction, err := h.transactionService.Create(c.Request.Context(), claims.UserID, &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.CreatedResponse(c, "Transaction created successfully", transaction)
}

// UpdateStatus handles PATCH /api/v1/transactions/:id/status
func (h *TransactionHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateTransactionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if errors, ok := utils.Validate(&req); !ok {
		utils.ValidationErrorResponse(c, errors)
		return
	}

	transaction, err := h.transactionService.UpdateStatus(c.Request.Context(), id, &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Transaction status updated successfully", transaction)
}
