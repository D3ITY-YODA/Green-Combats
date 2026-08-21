package httpx

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"green-compass-backend/pkg/clock"
)

func RequestLogger(logger *slog.Logger, clk clock.Clock) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := clk.Now()

		c.Next()

		logger.LogAttrs(c.Request.Context(), slog.LevelInfo, "http_request",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.Int64("duration_ms", clk.Now().Sub(start).Milliseconds()),
			slog.String("client_ip", c.ClientIP()),
		)
	}
}
