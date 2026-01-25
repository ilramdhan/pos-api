package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ilramdhan/pos-api/internal/service"
	"github.com/ilramdhan/pos-api/internal/utils"
)

// ReportHandler handles report endpoints
type ReportHandler struct {
	reportService *service.ReportService
}

// NewReportHandler creates a new report handler
func NewReportHandler(reportService *service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

// DailySales handles GET /api/v1/reports/sales/daily
func (h *ReportHandler) DailySales(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	// Default to last 30 days if not specified
	if dateFrom == "" {
		dateFrom = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	if dateTo == "" {
		dateTo = time.Now().Format("2006-01-02")
	}

	reports, err := h.reportService.GetDailySales(c.Request.Context(), dateFrom, dateTo)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Daily sales report retrieved successfully", reports)
}

// MonthlySales handles GET /api/v1/reports/sales/monthly
func (h *ReportHandler) MonthlySales(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	// Default to last 12 months if not specified
	if dateFrom == "" {
		dateFrom = time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
	}
	if dateTo == "" {
		dateTo = time.Now().Format("2006-01-02")
	}

	reports, err := h.reportService.GetMonthlySales(c.Request.Context(), dateFrom, dateTo)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Monthly sales report retrieved successfully", reports)
}

// TopProducts handles GET /api/v1/reports/products/top
func (h *ReportHandler) TopProducts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	// Default to last 30 days if not specified
	if dateFrom == "" {
		dateFrom = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	if dateTo == "" {
		dateTo = time.Now().Format("2006-01-02")
	}

	reports, err := h.reportService.GetTopProducts(c.Request.Context(), limit, dateFrom, dateTo)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Top products report retrieved successfully", reports)
}
