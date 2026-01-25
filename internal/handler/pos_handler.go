package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ilramdhan/pos-api/internal/dto"
	"github.com/ilramdhan/pos-api/internal/middleware"
	"github.com/ilramdhan/pos-api/internal/service"
	"github.com/ilramdhan/pos-api/internal/utils"
)

// POSHandler handles POS (Point of Sale) endpoints
type POSHandler struct {
	productService     *service.ProductService
	transactionService *service.TransactionService
}

// NewPOSHandler creates a new POS handler
func NewPOSHandler(productService *service.ProductService, transactionService *service.TransactionService) *POSHandler {
	return &POSHandler{
		productService:     productService,
		transactionService: transactionService,
	}
}

// GetProducts handles GET /api/v1/pos/products
func (h *POSHandler) GetProducts(c *gin.Context) {
	categoryID := c.Query("category_id")
	search := c.Query("search")

	products, _, err := h.productService.List(c.Request.Context(), dto.ProductListFilter{
		CategoryID: categoryID,
		Search:     search,
	}, utils.Pagination{Page: 1, PerPage: 100})
	if err != nil {
		// Log error for debugging
		fmt.Printf("[POS] Error fetching products: %v\n", err)
		utils.InternalServerError(c, "Failed to fetch products: "+err.Error())
		return
	}

	// Initialize as empty slice (not nil) to return [] instead of null in JSON
	posProducts := make([]gin.H, 0)
	for _, p := range products {
		categoryName := ""
		if p.Category != nil {
			categoryName = p.Category.Name
		}
		posProducts = append(posProducts, gin.H{
			"id":            p.ID,
			"name":          p.Name,
			"sku":           p.SKU,
			"price":         p.Price,
			"stock":         p.Stock,
			"category_id":   p.CategoryID,
			"category_name": categoryName,
			"image_url":     p.ImageURL,
		})
	}

	utils.SuccessResponse(c, http.StatusOK, "Products retrieved", posProducts)
}

// CreateTransaction handles POST /api/v1/pos/transactions
func (h *POSHandler) CreateTransaction(c *gin.Context) {
	var req struct {
		CustomerID    *string `json:"customer_id"`
		PaymentMethod string  `json:"payment_method"`
		Items         []struct {
			ProductID string  `json:"product_id"`
			Quantity  int     `json:"quantity"`
			UnitPrice float64 `json:"unit_price"`
			Discount  float64 `json:"discount"`
		} `json:"items"`
		Subtotal       float64 `json:"subtotal"`
		TaxAmount      float64 `json:"tax_amount"`
		DiscountAmount float64 `json:"discount_amount"`
		TotalAmount    float64 `json:"total_amount"`
		AmountPaid     float64 `json:"amount_paid"`
		ChangeAmount   float64 `json:"change_amount"`
		Notes          string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	claims := middleware.GetCurrentUser(c)
	if claims == nil {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	// Build DTO request
	txReq := &dto.CreateTransactionRequest{
		CustomerID:     req.CustomerID,
		PaymentMethod:  req.PaymentMethod,
		DiscountAmount: req.DiscountAmount,
		Notes:          req.Notes,
	}

	for _, item := range req.Items {
		txReq.Items = append(txReq.Items, dto.CreateTransactionItemDTO{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	// Create transaction
	tx, err := h.transactionService.Create(c.Request.Context(), claims.UserID, txReq)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.CreatedResponse(c, "Transaction created successfully", gin.H{
		"id":             tx.ID,
		"invoice_number": tx.InvoiceNumber,
		"status":         tx.Status,
		"total_amount":   tx.TotalAmount,
		"created_at":     tx.CreatedAt.Format(time.RFC3339),
	})
}

// HoldTransaction represents a held transaction
type HoldTransaction struct {
	ID           string    `json:"id"`
	HoldNumber   string    `json:"hold_number"`
	CustomerName string    `json:"customer_name"`
	ItemsCount   int       `json:"items_count"`
	TotalAmount  float64   `json:"total_amount"`
	CreatedAt    time.Time `json:"created_at"`
}

// Mock storage for held transactions (in production, use database)
var heldTransactions = make(map[string]*HoldTransaction)
var holdCounter = 0

// GetHeldTransactions handles GET /api/v1/pos/hold
func (h *POSHandler) GetHeldTransactions(c *gin.Context) {
	var holds []*HoldTransaction
	for _, hold := range heldTransactions {
		holds = append(holds, hold)
	}

	utils.SuccessResponse(c, http.StatusOK, "Held transactions retrieved", holds)
}

// HoldTransactionCreate handles POST /api/v1/pos/hold
func (h *POSHandler) HoldTransactionCreate(c *gin.Context) {
	var req struct {
		CustomerID   *string `json:"customer_id"`
		CustomerName string  `json:"customer_name"`
		Items        []struct {
			ProductID string `json:"product_id"`
			Quantity  int    `json:"quantity"`
		} `json:"items"`
		TotalAmount float64 `json:"total_amount"`
		Notes       string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	holdCounter++
	hold := &HoldTransaction{
		ID:           uuid.New().String(),
		HoldNumber:   "HOLD-" + strconv.Itoa(holdCounter),
		CustomerName: req.CustomerName,
		ItemsCount:   len(req.Items),
		TotalAmount:  req.TotalAmount,
		CreatedAt:    time.Now(),
	}

	heldTransactions[hold.ID] = hold

	utils.CreatedResponse(c, "Transaction held successfully", hold)
}

// DeleteHeldTransaction handles DELETE /api/v1/pos/hold/:id
func (h *POSHandler) DeleteHeldTransaction(c *gin.Context) {
	id := c.Param("id")
	delete(heldTransactions, id)
	utils.SuccessResponse(c, http.StatusOK, "Held transaction deleted", nil)
}
