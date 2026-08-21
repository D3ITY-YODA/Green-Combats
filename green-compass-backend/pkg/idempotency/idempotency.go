package idempotency

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync"
	"time"

	"green-compass-backend/pkg/cache"
)

// Store tracks idempotency keys and their associated response/status.
type Store struct {
	cache  cache.Cache
	prefix string
}

// NewStore creates a new idempotency store backed by the given cache.
func NewStore(c cache.Cache) *Store {
	return &Store{
		cache:  c,
		prefix: "idempotency:",
	}
}

// Result represents the stored result of a previous request.
type Result struct {
	Key        string      `json:"key"`
	StatusCode int         `json:"status_code"`
	Body       interface{} `json:"body"`
	CreatedAt  time.Time   `json:"created_at"`
}

// ComputeKey generates an idempotency key from request components.
// Typically: method + path + user_id + body_hash.
func ComputeKey(method, path, userID string, body []byte) string {
	h := sha256.New()
	h.Write([]byte(method))
	h.Write([]byte(path))
	h.Write([]byte(userID))
	h.Write(body)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Check returns the stored result if the key has already been processed.
// Returns nil, nil if the key is new (no prior result found).
func (s *Store) Check(ctx context.Context, key string) (*Result, error) {
	var result Result
	err := s.cache.Get(ctx, s.prefix+key, &result)
	if err != nil {
		if err == cache.ErrCacheMiss {
			return nil, nil // key is new
		}
		return nil, fmt.Errorf("check idempotency: %w", err)
	}
	return &result, nil
}

// Store saves the result of a processed request.
func (s *Store) Store(ctx context.Context, key string, statusCode int, body interface{}, ttl time.Duration) error {
	result := Result{
		Key:        key,
		StatusCode: statusCode,
		Body:       body,
		CreatedAt:  time.Now(),
	}
	return s.cache.Set(ctx, s.prefix+key, result, ttl)
}

// Invalidate removes an idempotency key (e.g., on error retry).
func (s *Store) Invalidate(ctx context.Context, key string) error {
	return s.cache.Delete(ctx, s.prefix+key)
}

// --- In-memory idempotency store (simpler, no cache dependency) ---

// MemoryStore is a standalone in-memory idempotency store.
type MemoryStore struct {
	mu    sync.RWMutex
	items map[string]Result
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		items: make(map[string]Result),
	}
}

func (m *MemoryStore) Check(_ context.Context, key string) (*Result, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result, ok := m.items[key]
	if !ok {
		return nil, nil
	}
	return &result, nil
}

func (m *MemoryStore) Store(_ context.Context, key string, statusCode int, body interface{}, _ time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items[key] = Result{
		Key:        key,
		StatusCode: statusCode,
		Body:       body,
		CreatedAt:  time.Now(),
	}
	return nil
}

func (m *MemoryStore) Invalidate(_ context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.items, key)
	return nil
}
