package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ilramdhan/pos-api/internal/utils"
	"golang.org/x/time/rate"
)

// IPRateLimiter holds rate limiters per IP address
type IPRateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rps      rate.Limit
	burst    int
}

// NewIPRateLimiter creates a new IP-based rate limiter
func NewIPRateLimiter(rps int, burst int) *IPRateLimiter {
	return &IPRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rps:      rate.Limit(rps),
		burst:    burst,
	}
}

// getLimiter returns the rate limiter for the given IP address
func (i *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	i.mu.RLock()
	limiter, exists := i.limiters[ip]
	i.mu.RUnlock()

	if exists {
		return limiter
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	// Double-check after acquiring write lock
	if limiter, exists := i.limiters[ip]; exists {
		return limiter
	}

	limiter = rate.NewLimiter(i.rps, i.burst)
	i.limiters[ip] = limiter
	return limiter
}

// Cleanup removes old limiters to prevent memory leaks
func (i *IPRateLimiter) Cleanup(maxAge time.Duration) {
	i.mu.Lock()
	defer i.mu.Unlock()

	// Simple cleanup: if map is too large, clear it
	if len(i.limiters) > 10000 {
		i.limiters = make(map[string]*rate.Limiter)
	}
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		l := limiter.getLimiter(ip)

		if !l.Allow() {
			utils.TooManyRequests(c)
			c.Abort()
			return
		}

		c.Next()
	}
}

// StartCleanup starts a background goroutine to clean up old limiters
func (i *IPRateLimiter) StartCleanup(interval, maxAge time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			i.Cleanup(maxAge)
		}
	}()
}
