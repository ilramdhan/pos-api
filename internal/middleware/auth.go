package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ilramdhan/pos-api/internal/utils"
)

const (
	// AuthorizationHeader is the header key for authorization
	AuthorizationHeader = "Authorization"
	// BearerPrefix is the prefix for bearer tokens
	BearerPrefix = "Bearer "
	// UserContextKey is the key for storing user claims in context
	UserContextKey = "user"
)

// AuthMiddleware creates a JWT authentication middleware
// It checks both Authorization header and httpOnly cookie
func AuthMiddleware(jwtManager *utils.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// First try Authorization header
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader != "" && strings.HasPrefix(authHeader, BearerPrefix) {
			tokenString = strings.TrimPrefix(authHeader, BearerPrefix)
		}

		// If no header token, try cookie
		if tokenString == "" {
			cookieToken, err := utils.GetAccessTokenFromCookie(c)
			if err == nil && cookieToken != "" {
				tokenString = cookieToken
			}
		}

		// No token found
		if tokenString == "" {
			utils.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			utils.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Store claims in context
		c.Set(UserContextKey, claims)
		c.Next()
	}
}

// RequireRole creates a middleware that requires specific roles
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get(UserContextKey)
		if !exists {
			utils.Unauthorized(c, "User not authenticated")
			c.Abort()
			return
		}

		userClaims, ok := claims.(*utils.JWTClaims)
		if !ok {
			utils.Unauthorized(c, "Invalid user claims")
			c.Abort()
			return
		}

		// Check if user role is in the allowed roles
		allowed := false
		for _, role := range roles {
			if userClaims.Role == role {
				allowed = true
				break
			}
		}

		if !allowed {
			utils.Forbidden(c, "You don't have permission to access this resource")
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetCurrentUser retrieves the current user claims from context
func GetCurrentUser(c *gin.Context) *utils.JWTClaims {
	claims, exists := c.Get(UserContextKey)
	if !exists {
		return nil
	}

	userClaims, ok := claims.(*utils.JWTClaims)
	if !ok {
		return nil
	}

	return userClaims
}
