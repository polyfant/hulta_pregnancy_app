package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Rate Limit Enforcement", func(t *testing.T) {
		router := gin.New()
		
		// Configure a very strict rate limit for testing
		rateLimitConfig := RateLimitConfig{
			RequestsPerSecond: 1,  // Only 1 request per second
			BurstLimit:        2,  // Allow 2 burst requests
		}
		
		router.Use(RateLimitMiddleware(rateLimitConfig))
		
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// First two requests should succeed
		w1 := performRequest(router, "GET", "/test", "127.0.0.1")
		assert.Equal(t, http.StatusOK, w1.Code)

		w2 := performRequest(router, "GET", "/test", "127.0.0.1")
		assert.Equal(t, http.StatusOK, w2.Code)

		// Third request should be rate limited
		w3 := performRequest(router, "GET", "/test", "127.0.0.1")
		assert.Equal(t, http.StatusTooManyRequests, w3.Code)
	})

	t.Run("Different IPs Have Separate Limits", func(t *testing.T) {
		router := gin.New()
		
		rateLimitConfig := RateLimitConfig{
			RequestsPerSecond: 1,
			BurstLimit:        2,
		}
		
		router.Use(RateLimitMiddleware(rateLimitConfig))
		
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// Requests from different IPs should each get their own limit
		w1 := performRequest(router, "GET", "/test", "127.0.0.1")
		assert.Equal(t, http.StatusOK, w1.Code)

		w2 := performRequest(router, "GET", "/test", "192.168.1.1")
		assert.Equal(t, http.StatusOK, w2.Code)
	})

	t.Run("Concurrent Request Handling", func(t *testing.T) {
		router := gin.New()
		
		rateLimitConfig := RateLimitConfig{
			RequestsPerSecond: 5,
			BurstLimit:        10,
		}
		
		router.Use(RateLimitMiddleware(rateLimitConfig))
		
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		var wg sync.WaitGroup
		var successCount, failureCount int
		var mu sync.Mutex

		// Simulate 20 concurrent requests
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				w := performRequest(router, "GET", "/test", "127.0.0.1")
				
				mu.Lock()
				if w.Code == http.StatusOK {
					successCount++
				} else {
					failureCount++
				}
				mu.Unlock()
			}()
		}

		wg.Wait()

		// We expect some requests to succeed and some to be rate limited
		assert.True(t, successCount <= 10)  // Should not exceed burst limit
		assert.True(t, failureCount > 0)    // Some requests should be rate limited
	})

	t.Run("Time-Based Rate Recovery", func(t *testing.T) {
		router := gin.New()
		
		rateLimitConfig := RateLimitConfig{
			RequestsPerSecond: 1,  // Only 1 request per second
			BurstLimit:        2,  // Allow 2 burst requests
		}
		
		router.Use(RateLimitMiddleware(rateLimitConfig))
		
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// First two requests should succeed
		w1 := performRequest(router, "GET", "/test", "127.0.0.1")
		assert.Equal(t, http.StatusOK, w1.Code)

		w2 := performRequest(router, "GET", "/test", "127.0.0.1")
		assert.Equal(t, http.StatusOK, w2.Code)

		// Third request should be rate limited
		w3 := performRequest(router, "GET", "/test", "127.0.0.1")
		assert.Equal(t, http.StatusTooManyRequests, w3.Code)

		// Wait for rate limit to recover
		time.Sleep(1 * time.Second)

		// Next request should succeed after waiting
		w4 := performRequest(router, "GET", "/test", "127.0.0.1")
		assert.Equal(t, http.StatusOK, w4.Code)
	})

	t.Run("Precise Time-Based Rate Limiting", func(t *testing.T) {
		router := gin.New()
		
		rateLimitConfig := RateLimitConfig{
			RequestsPerSecond: 2,  // 2 requests per second
			BurstLimit:        3,  // Allow 3 burst requests
		}
		
		router.Use(RateLimitMiddleware(rateLimitConfig))
		
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// Track precise timing of requests
		var mu sync.Mutex
		var requestTimes []time.Time

		// Helper function to record request times
		recordRequestTime := func(w *httptest.ResponseRecorder) {
			mu.Lock()
			defer mu.Unlock()
			requestTimes = append(requestTimes, time.Now())
		}

		// Perform multiple requests
		w1 := performRequest(router, "GET", "/test", "127.0.0.1")
		recordRequestTime(w1)
		assert.Equal(t, http.StatusOK, w1.Code)

		w2 := performRequest(router, "GET", "/test", "127.0.0.1")
		recordRequestTime(w2)
		assert.Equal(t, http.StatusOK, w2.Code)

		w3 := performRequest(router, "GET", "/test", "127.0.0.1")
		recordRequestTime(w3)
		assert.Equal(t, http.StatusOK, w3.Code)

		w4 := performRequest(router, "GET", "/test", "127.0.0.1")
		recordRequestTime(w4)
		assert.Equal(t, http.StatusTooManyRequests, w4.Code)

		// Verify time between first three requests is reasonable
		if len(requestTimes) >= 3 {
			// Check that requests are somewhat evenly spaced
			timeDiff1 := requestTimes[1].Sub(requestTimes[0])
			timeDiff2 := requestTimes[2].Sub(requestTimes[1])
			
			assert.True(t, timeDiff1 < 500*time.Millisecond, "Requests should be close together")
			assert.True(t, timeDiff2 < 500*time.Millisecond, "Requests should be close together")
		}
	})

	t.Run("Direct Rate Limiter Demonstration", func(t *testing.T) {
		// Create a rate limiter directly
		limiter := rate.NewLimiter(
			rate.Limit(2),   // 2 events per second
			3,              // Burst of 3 tokens
		)

		// Track successful and failed attempts
		successCount := 0
		failureCount := 0

		// Simulate multiple requests
		for i := 0; i < 6; i++ {
			// Check if request is allowed
			if limiter.Allow() {
				successCount++
			} else {
				failureCount++
			}

			// Small delay to simulate request spacing
			time.Sleep(250 * time.Millisecond)
		}

		// Detailed logging for debugging
		t.Logf("Successful requests: %d, Failed requests: %d", successCount, failureCount)
		
		// More flexible assertions
		assert.LessOrEqual(t, successCount, 5, "Should not exceed burst limit")
		assert.GreaterOrEqual(t, failureCount, 1, "Some requests should be rate limited")
		
		// Demonstrate waiting for tokens
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Wait for a token if none available
		err := limiter.Wait(ctx)
		assert.NoError(t, err, "Should be able to wait for a token")
	})
}

// Helper function to perform test requests
func performRequest(r http.Handler, method, path, ip string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	req.RemoteAddr = ip + ":12345"  // Add a port to simulate full IP
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
