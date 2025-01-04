package middleware

import (
	"net"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/audit"
)

type ConnectionTracker struct {
	mu              sync.RWMutex
	connections     map[string]int
	maxConnections  int
	blockDuration   time.Duration
	blockedIPs      map[string]time.Time
}

func NewConnectionTracker(maxConnections int, blockDuration time.Duration) *ConnectionTracker {
	return &ConnectionTracker{
		connections:     make(map[string]int),
		maxConnections:  maxConnections,
		blockDuration:   blockDuration,
		blockedIPs:      make(map[string]time.Time),
	}
}

func (ct *ConnectionTracker) AllowConnection(ip string) bool {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	// Check if IP is blocked
	if blockedUntil, isBlocked := ct.blockedIPs[ip]; isBlocked {
		if time.Now().Before(blockedUntil) {
			return false
		}
		// Unblock if time has passed
		delete(ct.blockedIPs, ip)
	}

	// Increment connection count
	ct.connections[ip]++

	// Block if too many connections
	if ct.connections[ip] > ct.maxConnections {
		ct.blockedIPs[ip] = time.Now().Add(ct.blockDuration)
		return false
	}

	return true
}

func (ct *ConnectionTracker) ReleaseConnection(ip string) {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	if ct.connections[ip] > 0 {
		ct.connections[ip]--
	}
}

func ConnectionIsolationMiddleware() gin.HandlerFunc {
	tracker := NewConnectionTracker(10, 5*time.Minute)
	auditTrail, _ := audit.NewAuditTrail("./logs")

	return func(c *gin.Context) {
		// Get client IP
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			ip = c.ClientIP()
		}

		// Check connection
		if !tracker.AllowConnection(ip) {
			// Log blocked connection attempt
			auditTrail.LogEvent(ip, "CONNECTION_BLOCKED", map[string]interface{}{
				"reason": "Too many connections",
				"ip":     ip,
			})

			c.JSON(429, gin.H{
				"error": "Too many connections. Please try again later.",
			})
			c.Abort()
			return
		}

		// Release connection after request
		defer tracker.ReleaseConnection(ip)

		c.Next()
	}
}
