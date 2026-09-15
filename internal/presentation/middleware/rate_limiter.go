package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type clientBucket struct {
	tokens     float64
	lastUpdate time.Time
}

type IPRateLimiter struct {
	mu          sync.Mutex
	clients     map[string]*clientBucket
	rate        float64 // tokens added per second
	burst       float64 // max tokens
	cleanupFreq time.Duration
}

func NewIPRateLimiter(ratePerMinute float64, burst float64) *IPRateLimiter {
	limiter := &IPRateLimiter{
		clients:     make(map[string]*clientBucket),
		rate:        ratePerMinute / 60.0,
		burst:       burst,
		cleanupFreq: 10 * time.Minute,
	}

	go limiter.cleanupLoop()

	return limiter
}

func (l *IPRateLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	bucket, exists := l.clients[ip]
	if !exists {
		l.clients[ip] = &clientBucket{
			tokens:     l.burst - 1,
			lastUpdate: now,
		}
		return true
	}

	// Refill tokens based on elapsed time
	elapsed := now.Sub(bucket.lastUpdate).Seconds()
	bucket.tokens += elapsed * l.rate
	if bucket.tokens > l.burst {
		bucket.tokens = l.burst
	}
	bucket.lastUpdate = now

	if bucket.tokens >= 1 {
		bucket.tokens -= 1
		return true
	}

	return false
}

func (l *IPRateLimiter) cleanupLoop() {
	ticker := time.NewTicker(l.cleanupFreq)
	for range ticker.C {
		l.mu.Lock()
		now := time.Now()
		for ip, bucket := range l.clients {
			// Remove entries inactive for more than 15 minutes
			if now.Sub(bucket.lastUpdate) > 15*time.Minute {
				delete(l.clients, ip)
			}
		}
		l.mu.Unlock()
	}
}

// RateLimitMiddleware creates Gin middleware using the given IPRateLimiter.
func RateLimitMiddleware(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		if !limiter.allow(clientIP) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests, please try again later",
			})
			return
		}
		c.Next()
	}
}
