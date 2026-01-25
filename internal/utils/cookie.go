package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// AccessTokenCookieName is the cookie name for access token
	AccessTokenCookieName = "access_token"
	// RefreshTokenCookieName is the cookie name for refresh token
	RefreshTokenCookieName = "refresh_token"
)

// CookieConfig holds cookie configuration
type CookieConfig struct {
	Domain   string
	Secure   bool
	SameSite http.SameSite
}

// DefaultCookieConfig returns default cookie configuration
func DefaultCookieConfig(isProduction bool) CookieConfig {
	return CookieConfig{
		Domain:   "",
		Secure:   isProduction,
		SameSite: http.SameSiteLaxMode,
	}
}

// SetAuthCookies sets access and refresh tokens as httpOnly cookies
func SetAuthCookies(c *gin.Context, accessToken, refreshToken string, accessExpiry, refreshExpiry time.Duration, cfg CookieConfig) {
	// Access token cookie
	c.SetSameSite(cfg.SameSite)
	c.SetCookie(
		AccessTokenCookieName,
		accessToken,
		int(accessExpiry.Seconds()),
		"/",
		cfg.Domain,
		cfg.Secure,
		true, // HttpOnly
	)

	// Refresh token cookie
	c.SetCookie(
		RefreshTokenCookieName,
		refreshToken,
		int(refreshExpiry.Seconds()),
		"/api/v1/auth", // Only send to auth endpoints
		cfg.Domain,
		cfg.Secure,
		true, // HttpOnly
	)
}

// ClearAuthCookies clears authentication cookies (for logout)
func ClearAuthCookies(c *gin.Context, cfg CookieConfig) {
	c.SetSameSite(cfg.SameSite)

	// Clear access token
	c.SetCookie(
		AccessTokenCookieName,
		"",
		-1,
		"/",
		cfg.Domain,
		cfg.Secure,
		true,
	)

	// Clear refresh token
	c.SetCookie(
		RefreshTokenCookieName,
		"",
		-1,
		"/api/v1/auth",
		cfg.Domain,
		cfg.Secure,
		true,
	)
}

// GetAccessTokenFromCookie retrieves access token from cookie
func GetAccessTokenFromCookie(c *gin.Context) (string, error) {
	return c.Cookie(AccessTokenCookieName)
}

// GetRefreshTokenFromCookie retrieves refresh token from cookie
func GetRefreshTokenFromCookie(c *gin.Context) (string, error) {
	return c.Cookie(RefreshTokenCookieName)
}
