package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ilramdhan/pos-api/internal/dto"
	"github.com/ilramdhan/pos-api/internal/service"
	"github.com/ilramdhan/pos-api/internal/utils"
)

// DashboardHandler handles dashboard and statistics endpoints
type DashboardHandler struct {
	transactionService *service.TransactionService
	productService     *service.ProductService
	categoryService    *service.CategoryService
	customerService    *service.CustomerService
	reportService      *service.ReportService
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(
	transactionService *service.TransactionService,
	productService *service.ProductService,
	categoryService *service.CategoryService,
	customerService *service.CustomerService,
	reportService *service.ReportService,
) *DashboardHandler {
	return &DashboardHandler{
		transactionService: transactionService,
		productService:     productService,
		categoryService:    categoryService,
		customerService:    customerService,
		reportService:      reportService,
	}
}

// GetDashboardStats handles GET /api/v1/reports/dashboard/stats
func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	ctx := c.Request.Context()
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	// Get today's sales
	todaySales, _ := h.reportService.GetDailySales(ctx, today, today)
	yesterdaySales, _ := h.reportService.GetDailySales(ctx, yesterday, yesterday)

	var todayAmount float64
	var yesterdayAmount float64
	var todayTransactions int

	if len(todaySales) > 0 {
		todayAmount = todaySales[0].TotalAmount
		todayTransactions = todaySales[0].TotalTransactions
	}
	if len(yesterdaySales) > 0 {
		yesterdayAmount = yesterdaySales[0].TotalAmount
	}

	changePercent := 0.0
	if yesterdayAmount > 0 {
		changePercent = ((todayAmount - yesterdayAmount) / yesterdayAmount) * 100
	}

	utils.SuccessResponse(c, http.StatusOK, "Dashboard stats retrieved", gin.H{
		"today_sales": gin.H{
			"amount":         todayAmount,
			"change_percent": changePercent,
			"comparison":     "yesterday",
		},
		"active_orders": gin.H{
			"count":              todayTransactions,
			"avg_prep_time_mins": 12,
		},
		"net_margin": gin.H{
			"percent":        32.4,
			"change_percent": 2.1,
			"comparison":     "last_week",
		},
	})
}

// GetRealtimeSales handles GET /api/v1/reports/sales/realtime
func (h *DashboardHandler) GetRealtimeSales(c *gin.Context) {
	interval := c.DefaultQuery("interval", "hourly")
	date := c.DefaultQuery("date", time.Now().Format("2006-01-02"))

	// Generate mock data points based on interval
	var dataPoints []gin.H
	if interval == "hourly" {
		for hour := 8; hour <= time.Now().Hour() && hour <= 22; hour++ {
			dataPoints = append(dataPoints, gin.H{
				"time":   time.Date(2006, 1, 2, hour, 0, 0, 0, time.Local).Format("15:04"),
				"amount": float64(100 + (hour * 50)),
			})
		}
	} else {
		for i := 0; i < 24; i++ {
			mins := i * 15
			dataPoints = append(dataPoints, gin.H{
				"time":   time.Date(2006, 1, 2, 8+mins/60, mins%60, 0, 0, time.Local).Format("15:04"),
				"amount": float64(50 + (i * 20)),
			})
		}
	}

	utils.SuccessResponse(c, http.StatusOK, "Realtime sales data retrieved", gin.H{
		"interval":    interval,
		"date":        date,
		"data_points": dataPoints,
	})
}

// GetRecentTransactions handles GET /api/v1/transactions/recent
func (h *DashboardHandler) GetRecentTransactions(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "5")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 5
	}
	if limit > 20 {
		limit = 20
	}

	pagination := utils.Pagination{Page: 1, PerPage: limit, Sort: "created_at", Order: "desc"}
	transactions, _, err := h.transactionService.List(c.Request.Context(),
		dto.TransactionListFilter{}, pagination)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	var recentTx []gin.H
	for _, tx := range transactions {
		tableOrType := "Walk-in"
		if tx.CustomerID != nil {
			tableOrType = "Member"
		}

		recentTx = append(recentTx, gin.H{
			"id":            tx.ID,
			"order_number":  tx.InvoiceNumber,
			"table_or_type": tableOrType,
			"total_amount":  tx.TotalAmount,
			"status":        tx.Status,
			"created_at":    tx.CreatedAt.Format(time.RFC3339),
			"time_ago":      timeAgo(tx.CreatedAt),
		})
	}

	utils.SuccessResponse(c, http.StatusOK, "Recent transactions retrieved", recentTx)
}

// GetSystemHealth handles GET /api/v1/system/health/detailed
func (h *DashboardHandler) GetSystemHealth(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "System health retrieved", gin.H{
		"terminals": []gin.H{
			{"id": "POS-01", "type": "pos", "latency_ms": 35, "status": "online"},
			{"id": "POS-02", "type": "pos", "latency_ms": 42, "status": "online"},
			{"id": "KDS-01", "type": "kds", "latency_ms": 32, "status": "online"},
		},
		"servers": []gin.H{
			{"id": "Server A", "uptime_percent": 99.9, "status": "online"},
			{"id": "Server B", "uptime_percent": 99.5, "status": "online"},
		},
		"integrations": []gin.H{
			{"name": "Payment Gateway", "status": "ok"},
			{"name": "Email Service", "status": "ok"},
			{"name": "SMS Gateway", "status": "ok"},
		},
	})
}

// GetCustomerStats handles GET /api/v1/customers/stats
func (h *DashboardHandler) GetCustomerStats(c *gin.Context) {
	pagination := utils.Pagination{Page: 1, PerPage: 1}
	_, total, _ := h.customerService.List(c.Request.Context(),
		dto.CustomerListFilter{}, pagination)

	utils.SuccessResponse(c, http.StatusOK, "Customer stats retrieved", gin.H{
		"total_customers":    total,
		"change_percent":     3.2,
		"comparison":         "last_month",
		"new_this_month":     int(float64(total) * 0.03),
		"new_change_percent": 12.8,
		"avg_loyalty_points": 150,
	})
}

// GetProductStats handles GET /api/v1/products/stats
func (h *DashboardHandler) GetProductStats(c *gin.Context) {
	pagination := utils.Pagination{Page: 1, PerPage: 1}
	_, total, _ := h.productService.List(c.Request.Context(),
		dto.ProductListFilter{}, pagination)

	utils.SuccessResponse(c, http.StatusOK, "Product stats retrieved", gin.H{
		"total_sku":            total,
		"low_stock_count":      3,
		"out_of_stock_count":   1,
		"inventory_value":      15000000,
		"value_change_percent": 2.4,
		"comparison":           "last_week",
	})
}

// GetStockMovements handles GET /api/v1/products/stock-movements
func (h *DashboardHandler) GetStockMovements(c *gin.Context) {
	// Return mock data for stock movements
	movements := []gin.H{
		{
			"id":              "1",
			"product_id":      "p1",
			"product_name":    "Kopi Susu Gula Aren",
			"sku":             "BEV-002",
			"type":            "sale",
			"quantity_change": -2,
			"new_balance":     148,
			"user":            "Cashier",
			"created_at":      time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
		},
		{
			"id":              "2",
			"product_id":      "p2",
			"product_name":    "Nasi Goreng Spesial",
			"sku":             "FOOD-001",
			"type":            "restock",
			"quantity_change": 50,
			"new_balance":     150,
			"user":            "Manager",
			"created_at":      time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
		},
	}

	utils.SuccessWithMeta(c, "Stock movements retrieved", movements, &utils.Meta{
		Page:       1,
		PerPage:    10,
		Total:      2,
		TotalPages: 1,
	})
}

// GetCategoryStats handles GET /api/v1/categories/stats
func (h *DashboardHandler) GetCategoryStats(c *gin.Context) {
	pagination := utils.Pagination{Page: 1, PerPage: 1}
	_, total, _ := h.categoryService.List(c.Request.Context(), pagination)

	utils.SuccessResponse(c, http.StatusOK, "Category stats retrieved", gin.H{
		"total_categories":  total,
		"new_this_week":     0,
		"active_categories": total,
		"coverage_percent":  100,
		"most_popular": gin.H{
			"id":         "cat1",
			"name":       "Beverages",
			"items_sold": 250,
		},
	})
}

// GetCategoryActivityLog handles GET /api/v1/categories/activity-log
func (h *DashboardHandler) GetCategoryActivityLog(c *gin.Context) {
	activities := []gin.H{
		{
			"id":         "1",
			"event":      "Created \"Personal Care\"",
			"user":       "Admin",
			"details":    "New category added",
			"created_at": time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
		},
	}

	utils.SuccessResponse(c, http.StatusOK, "Category activity log retrieved", activities)
}

// GetTransactionStats handles GET /api/v1/transactions/stats
func (h *DashboardHandler) GetTransactionStats(c *gin.Context) {
	ctx := c.Request.Context()
	today := time.Now().Format("2006-01-02")
	monthStart := time.Now().Format("2006-01") + "-01"

	// Get today's transactions
	todaySales, _ := h.reportService.GetDailySales(ctx, today, today)
	monthSales, _ := h.reportService.GetDailySales(ctx, monthStart, today)

	var todayTx, todayRevenue int
	var totalTx int
	var totalRevenue float64

	if len(todaySales) > 0 {
		todayTx = todaySales[0].TotalTransactions
		todayRevenue = int(todaySales[0].TotalAmount)
	}

	for _, day := range monthSales {
		totalTx += day.TotalTransactions
		totalRevenue += day.TotalAmount
	}

	avgOrderValue := 0.0
	if totalTx > 0 {
		avgOrderValue = totalRevenue / float64(totalTx)
	}

	utils.SuccessResponse(c, http.StatusOK, "Transaction stats retrieved", gin.H{
		"total_transactions": totalTx,
		"total_revenue":      totalRevenue,
		"avg_order_value":    avgOrderValue,
		"today_transactions": todayTx,
		"today_revenue":      todayRevenue,
	})
}

// GetWeeklySales handles GET /api/v1/reports/sales/weekly
func (h *DashboardHandler) GetWeeklySales(c *gin.Context) {
	dateFrom := c.DefaultQuery("date_from", time.Now().AddDate(0, 0, -28).Format("2006-01-02"))
	dateTo := c.DefaultQuery("date_to", time.Now().Format("2006-01-02"))

	dailySales, _ := h.reportService.GetDailySales(c.Request.Context(), dateFrom, dateTo)

	// Aggregate by week
	weeklyData := make(map[string]*struct {
		WeekStart    string
		Revenue      float64
		Transactions int
	})

	var totalRevenue float64
	var totalTx int

	for _, day := range dailySales {
		weekNum := getWeekNumber(day.Date)
		if _, ok := weeklyData[weekNum]; !ok {
			weeklyData[weekNum] = &struct {
				WeekStart    string
				Revenue      float64
				Transactions int
			}{
				WeekStart: day.Date,
			}
		}
		weeklyData[weekNum].Revenue += day.TotalAmount
		weeklyData[weekNum].Transactions += day.TotalTransactions
		totalRevenue += day.TotalAmount
		totalTx += day.TotalTransactions
	}

	var chartData []gin.H
	for week, data := range weeklyData {
		chartData = append(chartData, gin.H{
			"week":         week,
			"week_start":   data.WeekStart,
			"revenue":      data.Revenue,
			"transactions": data.Transactions,
		})
	}

	avgOrderValue := 0.0
	if totalTx > 0 {
		avgOrderValue = totalRevenue / float64(totalTx)
	}

	utils.SuccessResponse(c, http.StatusOK, "Weekly sales retrieved", gin.H{
		"period":             "weekly",
		"total_revenue":      totalRevenue,
		"total_transactions": totalTx,
		"avg_order_value":    avgOrderValue,
		"chart_data":         chartData,
	})
}

// GetCategoryPerformance handles GET /api/v1/reports/categories/performance
func (h *DashboardHandler) GetCategoryPerformance(c *gin.Context) {
	categories, _, _ := h.categoryService.List(c.Request.Context(), utils.Pagination{Page: 1, PerPage: 100})

	var performance []gin.H
	totalRevenue := 100000.0 // Mock total

	for i, cat := range categories {
		sold := 100 - (i * 20)
		revenue := float64(sold) * 25000
		percentage := (revenue / totalRevenue) * 100

		performance = append(performance, gin.H{
			"category_id":   cat.ID,
			"category_name": cat.Name,
			"total_sold":    sold,
			"total_revenue": revenue,
			"percentage":    percentage,
		})
	}

	utils.SuccessResponse(c, http.StatusOK, "Category performance retrieved", performance)
}

// Helper function to calculate time ago
func timeAgo(t time.Time) string {
	diff := time.Since(t)

	switch {
	case diff < time.Minute:
		return "Just now"
	case diff < time.Hour:
		mins := int(diff.Minutes())
		return strconv.Itoa(mins) + " mins ago"
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		return strconv.Itoa(hours) + " hours ago"
	default:
		days := int(diff.Hours() / 24)
		return strconv.Itoa(days) + " days ago"
	}
}

// Helper function to get ISO week number
func getWeekNumber(dateStr string) string {
	t, _ := time.Parse("2006-01-02", dateStr)
	year, week := t.ISOWeek()
	return strconv.Itoa(year) + "-W" + strconv.Itoa(week)
}
