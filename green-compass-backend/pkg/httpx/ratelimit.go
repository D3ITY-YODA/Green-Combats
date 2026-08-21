package httpx

import (
	"math"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter is a token-bucket rate limiter keyed by client identity.
type RateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	rate    float64   // tokens per second
	burst   float64   // bucket capacity
	now     func() time.Time
}

type bucket struct {
	tokens   float64
	lastSeen time.Time
}

// NewRateLimiter creates a limiter allowing `rate` requests per second with
// the given burst. Stale buckets are pruned lazily.
func NewRateLimiter(ratePerMinute, burst int) *RateLimiter {
	return &RateLimiter{
		buckets: make(map[string]*bucket),
		rate:    float64(ratePerMinute) / 60.0,
		burst:   float64(burst),
		now:     time.Now,
	}
}

// Allow reports whether the key may proceed and consumes one token if so.
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := rl.now()
	b, ok := rl.buckets[key]
	if !ok {
		rl.buckets[key] = &bucket{tokens: rl.burst - 1, lastSeen: now}
		return true
	}

	elapsed := now.Sub(b.lastSeen).Seconds()
	b.tokens = math.Min(rl.burst, b.tokens+elapsed*rl.rate)
	b.lastSeen = now

	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// Middleware returns a gin middleware enforcing per-client limits.
// Keys combine IP and (when present) authenticated user ID.
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		if id, ok := c.Get("request_id"); ok {
			_ = id // request_id is not identity; kept for clarity
		}
		if !rl.Allow(key) {
			c.Header("Retry-After", "30")
			HandleError(c, ErrTooManyRequests)
			c.Abort()
			return
		}
		c.Next()
	}
}

// Prune removes buckets idle longer than ttl. Call periodically from a
// background goroutine in long-lived processes.
func (rl *RateLimiter) Prune(ttl time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	cutoff := rl.now().Add(-ttl)
	for k, b := range rl.buckets {
		if b.lastSeen.Before(cutoff) {
			delete(rl.buckets, k)
		}
	}
}