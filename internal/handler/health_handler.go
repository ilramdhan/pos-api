package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilramdhan/pos-api/internal/config"
	"github.com/ilramdhan/pos-api/internal/utils"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	cfg *config.Config
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(cfg *config.Config) *HealthHandler {
	return &HealthHandler{cfg: cfg}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status      string `json:"status"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
}

// Check handles GET /health
func (h *HealthHandler) Check(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Service is healthy", HealthResponse{
		Status:      "ok",
		Version:     h.cfg.App.Version,
		Environment: h.cfg.App.Env,
	})
}

// CheckV1 handles GET /api/v1/health
func (h *HealthHandler) CheckV1(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Service is healthy", HealthResponse{
		Status:      "ok",
		Version:     h.cfg.App.Version,
		Environment: h.cfg.App.Env,
	})
}
