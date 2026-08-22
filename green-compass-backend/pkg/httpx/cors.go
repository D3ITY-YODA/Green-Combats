package httpx

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig configures cross-origin resource sharing for browser clients.
type CORSConfig struct {
	// AllowedOrigins is the list of origins permitted to access the API.
	// Use "*" to allow any origin (ignored when AllowCredentials is true).
	AllowedOrigins []string
	// AllowCredentials permits credentials such as cookies and authorization headers.
	AllowCredentials bool
}

// CORS returns a gin middleware that applies configurable CORS headers and
// short-circuits preflight (OPTIONS) requests.
func CORS(cfg CORSConfig) gin.HandlerFunc {
	allowed := normalizeOrigins(cfg.AllowedOrigins)
	allowAll := len(allowed) == 1 && allowed[0] == "*"

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			// Not a cross-origin request; nothing to do.
			c.Next()
			return
		}

		if allowAll {
			c.Header("Access-Control-Allow-Origin", "*")
		} else if matchesOrigin(allowed, origin) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept, X-Request-Id")

		if cfg.AllowCredentials && !allowAll {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func normalizeOrigins(in []string) []string {
	out := make([]string, 0, len(in))
	for _, o := range in {
		if o = strings.TrimSpace(o); o != "" {
			out = append(out, strings.ToLower(o))
		}
	}
	return out
}

func matchesOrigin(allowed []string, origin string) bool {
	origin = strings.ToLower(origin)
	for _, a := range allowed {
		if a == "*" || a == origin {
			return true
		}
	}
	return false
}
