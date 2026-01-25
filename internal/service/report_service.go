package service

import (
	"context"

	"github.com/ilramdhan/pos-api/internal/dto"
	"github.com/ilramdhan/pos-api/internal/repository"
)

// ReportService handles report generation
type ReportService struct {
	transactionRepo repository.TransactionRepository
}

// NewReportService creates a new report service
func NewReportService(transactionRepo repository.TransactionRepository) *ReportService {
	return &ReportService{
		transactionRepo: transactionRepo,
	}
}

// GetDailySales returns daily sales summary
func (s *ReportService) GetDailySales(ctx context.Context, dateFrom, dateTo string) ([]dto.DailySalesReport, error) {
	return s.transactionRepo.GetDailySales(ctx, dateFrom, dateTo)
}

// GetMonthlySales returns monthly sales summary
func (s *ReportService) GetMonthlySales(ctx context.Context, dateFrom, dateTo string) ([]dto.MonthlySalesReport, error) {
	return s.transactionRepo.GetMonthlySales(ctx, dateFrom, dateTo)
}

// GetTopProducts returns top selling products
func (s *ReportService) GetTopProducts(ctx context.Context, limit int, dateFrom, dateTo string) ([]dto.TopProductReport, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return s.transactionRepo.GetTopProducts(ctx, limit, dateFrom, dateTo)
}
