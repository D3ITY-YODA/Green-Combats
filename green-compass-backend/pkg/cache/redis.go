package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

// Cache provides key-value caching with optional TTL.
type Cache interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Close() error
}

// MemoryCache is an in-memory cache for development/testing.
type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]cacheEntry
}

type cacheEntry struct {
	value     []byte
	expiresAt time.Time
}

// NewMemoryCache creates an in-memory cache.
func NewMemoryCache() *MemoryCache {
	c := &MemoryCache{
		items: make(map[string]cacheEntry),
	}
	// Start background cleanup
	go c.cleanup()
	return c
}

func (c *MemoryCache) Get(_ context.Context, key string, dest interface{}) error {
	c.mu.RLock()
	entry, ok := c.items[key]
	c.mu.RUnlock()

	if !ok || time.Now().After(entry.expiresAt) {
		return ErrCacheMiss
	}

	if dest == nil {
		return nil
	}

	return json.Unmarshal(entry.value, dest)
}

func (c *MemoryCache) Set(_ context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal value: %w", err)
	}

	c.mu.Lock()
	c.items[key] = cacheEntry{
		value:     data,
		expiresAt: time.Now().Add(ttl),
	}
	c.mu.Unlock()
	return nil
}

func (c *MemoryCache) Delete(_ context.Context, key string) error {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
	return nil
}

func (c *MemoryCache) Exists(_ context.Context, key string) (bool, error) {
	c.mu.RLock()
	entry, ok := c.items[key]
	c.mu.RUnlock()

	if !ok {
		return false, nil
	}
	return time.Now().Before(entry.expiresAt), nil
}

func (c *MemoryCache) Close() error {
	return nil
}

func (c *MemoryCache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.items {
			if now.After(entry.expiresAt) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

// --- Errors ---

var ErrCacheMiss = errors.New("cache miss")

// --- Redis Cache (stub for when Redis is available) ---

// RedisCache wraps a Redis client for production use.
// When the Redis address is empty, it falls back to MemoryCache.
type RedisCache struct {
	*MemoryCache // fallback
	addr         string
}

// NewRedisCache creates a Redis-backed cache, falling back to memory if addr is empty.
func NewRedisCache(addr string) Cache {
	if addr == "" {
		return NewMemoryCache()
	}

	// In production, this would use go-redis client:
	// return &RedisCache{client: redis.NewClient(&redis.Options{Addr: addr})}
	// For now, fall back to memory cache
	return &RedisCache{
		MemoryCache: NewMemoryCache(),
		addr:        addr,
	}
}
