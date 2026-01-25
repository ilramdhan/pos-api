package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ilramdhan/pos-api/internal/config"
	"github.com/ilramdhan/pos-api/internal/dto"
	"github.com/ilramdhan/pos-api/internal/middleware"
	"github.com/ilramdhan/pos-api/internal/service"
	"github.com/ilramdhan/pos-api/internal/utils"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService   *service.AuthService
	cookieConfig  utils.CookieConfig
	jwtExpiry     time.Duration
	refreshExpiry time.Duration
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService:   authService,
		cookieConfig:  utils.DefaultCookieConfig(cfg.IsProduction()),
		jwtExpiry:     time.Duration(cfg.JWT.ExpiryHours) * time.Hour,
		refreshExpiry: time.Duration(cfg.JWT.RefreshExpiryHours) * time.Hour,
	}
}

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if errors, ok := utils.Validate(&req); !ok {
		utils.ValidationErrorResponse(c, errors)
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Set httpOnly cookies
	utils.SetAuthCookies(
		c,
		resp.Token.AccessToken,
		resp.Token.RefreshToken,
		h.jwtExpiry,
		h.refreshExpiry,
		h.cookieConfig,
	)

	// Return user info (tokens are in cookies, not in response for security)
	utils.SuccessResponse(c, http.StatusOK, "Login successful", gin.H{
		"user": resp.User,
		"token": gin.H{
			"expires_in": resp.Token.ExpiresIn,
			"token_type": resp.Token.TokenType,
		},
	})
}

// Register handles POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if errors, ok := utils.Validate(&req); !ok {
		utils.ValidationErrorResponse(c, errors)
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Set httpOnly cookies
	utils.SetAuthCookies(
		c,
		resp.Token.AccessToken,
		resp.Token.RefreshToken,
		h.jwtExpiry,
		h.refreshExpiry,
		h.cookieConfig,
	)

	utils.CreatedResponse(c, "Registration successful", gin.H{
		"user": resp.User,
		"token": gin.H{
			"expires_in": resp.Token.ExpiresIn,
			"token_type": resp.Token.TokenType,
		},
	})
}

// Me handles GET /api/v1/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	claims := middleware.GetCurrentUser(c)
	if claims == nil {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	resp, err := h.authService.GetCurrentUser(c.Request.Context(), claims.UserID)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User retrieved successfully", resp)
}

// UpdateProfile handles PUT /api/v1/auth/me
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	claims := middleware.GetCurrentUser(c)
	if claims == nil {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if errors, ok := utils.Validate(&req); !ok {
		utils.ValidationErrorResponse(c, errors)
		return
	}

	resp, err := h.authService.UpdateProfile(c.Request.Context(), claims.UserID, &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile updated successfully", resp)
}

// GetActivityLog handles GET /api/v1/auth/me/activity
func (h *AuthHandler) GetActivityLog(c *gin.Context) {
	claims := middleware.GetCurrentUser(c)
	if claims == nil {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	// Return mock activity data for now
	activities := []gin.H{
		{
			"id":         "1",
			"action":     "User Login",
			"device_ip":  c.ClientIP() + " (Browser)",
			"status":     "success",
			"created_at": time.Now().Format(time.RFC3339),
		},
	}

	utils.SuccessResponse(c, http.StatusOK, "Activity log retrieved", activities)
}

// Refresh handles POST /api/v1/auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	// Try to get refresh token from cookie first
	refreshToken, err := utils.GetRefreshTokenFromCookie(c)
	if err != nil || refreshToken == "" {
		// Fallback to request body
		var req dto.RefreshTokenRequest
		if err := c.ShouldBindJSON(&req); err == nil && req.RefreshToken != "" {
			refreshToken = req.RefreshToken
		}
	}

	if refreshToken == "" {
		utils.Unauthorized(c, "Refresh token required")
		return
	}

	resp, err := h.authService.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}

	// Set new cookies
	utils.SetAuthCookies(
		c,
		resp.AccessToken,
		resp.RefreshToken,
		h.jwtExpiry,
		h.refreshExpiry,
		h.cookieConfig,
	)

	utils.SuccessResponse(c, http.StatusOK, "Token refreshed successfully", gin.H{
		"expires_in": resp.ExpiresIn,
		"token_type": resp.TokenType,
	})
}

// Logout handles POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Clear cookies
	utils.ClearAuthCookies(c, h.cookieConfig)

	utils.SuccessResponse(c, http.StatusOK, "Successfully logged out", nil)
}

// ForgotPassword handles POST /api/v1/auth/forgot-password
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if errors, ok := utils.Validate(&req); !ok {
		utils.ValidationErrorResponse(c, errors)
		return
	}

	// In production, this would send an email
	// For now, we just return success (don't reveal if email exists)
	utils.SuccessResponse(c, http.StatusOK, "If the email exists, a password reset link has been sent", nil)
}

// ResetPassword handles POST /api/v1/auth/reset-password
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if errors, ok := utils.Validate(&req); !ok {
		utils.ValidationErrorResponse(c, errors)
		return
	}

	// In production, validate token and update password
	utils.SuccessResponse(c, http.StatusOK, "Password reset successful", nil)
}
