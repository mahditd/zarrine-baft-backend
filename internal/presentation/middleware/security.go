package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds baseline hardening headers (SRS 22).
// HTTPS itself is terminated at the production reverse proxy;
// HSTS is only sent when the request is already TLS to avoid
// breaking local HTTP development.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		if c.Request.TLS != nil {
			c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		c.Next()
	}
}
