package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestIPRateLimiter_BurstLimit(t *testing.T) {
	limiter := NewIPRateLimiter(60, 2) // 60/min, burst 2

	ip := "1.2.3.4"
	if !limiter.allow(ip) {
		t.Fatal("expected first request allowed")
	}
	if !limiter.allow(ip) {
		t.Fatal("expected second request allowed (burst)")
	}
	if limiter.allow(ip) {
		t.Fatal("expected third immediate request blocked")
	}
}

func TestIPRateLimiter_IsolatedPerIP(t *testing.T) {
	limiter := NewIPRateLimiter(60, 1)

	if !limiter.allow("10.0.0.1") {
		t.Fatal("expected first IP allowed")
	}
	if !limiter.allow("10.0.0.2") {
		t.Fatal("expected different IP allowed independently")
	}
	if limiter.allow("10.0.0.1") {
		t.Fatal("expected first IP now blocked")
	}
}

func TestRateLimitMiddleware_Returns429(t *testing.T) {
	gin.SetMode(gin.TestMode)
	limiter := NewIPRateLimiter(60, 1)

	router := gin.New()
	router.Use(RateLimitMiddleware(limiter))
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// First request passes
	req1 := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req1.RemoteAddr = "9.9.9.9:1234"
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	if w1.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w1.Code)
	}

	// Immediate second request from same IP is limited
	req2 := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req2.RemoteAddr = "9.9.9.9:1234"
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	if w2.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", w2.Code)
	}
}
