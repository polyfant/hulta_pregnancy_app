package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimitConfig defines the configuration for rate limiting
type RateLimitConfig struct {
	RequestsPerSecond float64
	BurstLimit        int
}

// RateLimitMiddleware creates a middleware to limit request rate
func RateLimitMiddleware(config RateLimitConfig) gin.HandlerFunc {
	// Create a map to store rate limiters per IP
	limiters := make(map[string]*rate.Limiter)
	var mu sync.RWMutex

	return func(c *gin.Context) {
		// Get client IP
		ip := c.ClientIP()

		mu.RLock()
		limiter, exists := limiters[ip]
		mu.RUnlock()

		if !exists {
			mu.Lock()
			// Create new rate limiter if not exists
			limiter = rate.NewLimiter(rate.Limit(config.RequestsPerSecond), config.BurstLimit)
			limiters[ip] = limiter
			mu.Unlock()
		}

		// Check if request is allowed
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please slow down.",
				"retry_after": time.Until(time.Now().Add(time.Second / 
					time.Duration(config.RequestsPerSecond))).Seconds(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
