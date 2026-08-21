package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSOptions configures the CORS middleware.
type CORSOptions struct {
	// AllowedOrigins is a list of origins allowed to make requests.
	// Use ["*"] to allow any origin (development only).
	AllowedOrigins []string
}

// CORSMiddleware returns a gin.HandlerFunc that sets CORS headers.
func CORSMiddleware(opts CORSOptions) gin.HandlerFunc {
	allowAll := false
	for _, o := range opts.AllowedOrigins {
		if o == "*" {
			allowAll = true
			break
		}
	}

	allowedSet := make(map[string]struct{}, len(opts.AllowedOrigins))
	for _, o := range opts.AllowedOrigins {
		allowedSet[o] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			c.Next()
			return
		}

		allowed := false
		if allowAll {
			allowed = true
		} else if _, ok := allowedSet[origin]; ok {
			allowed = true
		}

		if !allowed {
			c.Next()
			return
		}

		if allowAll {
			c.Header("Access-Control-Allow-Origin", "*")
		} else {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID, Accept, Accept-Language")
		c.Header("Access-Control-Expose-Headers", "X-Request-ID, Content-Length, Content-Type")
		c.Header("Access-Control-Max-Age", "86400")

		// Handle preflight
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
