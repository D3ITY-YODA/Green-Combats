package idempotency

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"green-compass-backend/pkg/cache"
)

func TestComputeKey_Deterministic(t *testing.T) {
	key1 := ComputeKey("POST", "/v1/users", "user123", []byte(`{"name":"test"}`))
	key2 := ComputeKey("POST", "/v1/users", "user123", []byte(`{"name":"test"}`))

	if key1 != key2 {
		t.Errorf("same inputs should produce same key: %s != %s", key1, key2)
	}
}

func TestComputeKey_DifferentInputs(t *testing.T) {
	key1 := ComputeKey("POST", "/v1/users", "user1", []byte(`{}`))
	key2 := ComputeKey("POST", "/v1/users", "user2", []byte(`{}`))

	if key1 == key2 {
		t.Error("different inputs should produce different keys")
	}
}

func TestMemoryStore_CheckMiss(t *testing.T) {
	store := NewMemoryStore()
	result, err := store.Check(context.Background(), "new-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Error("expected nil result for new key")
	}
}

func TestMemoryStore_StoreAndCheck(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	err := store.Store(ctx, "key1", 200, map[string]string{"status": "ok"}, 10*time.Minute)
	if err != nil {
		t.Fatalf("store failed: %v", err)
	}

	result, err := store.Check(ctx, "key1")
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.StatusCode != 200 {
		t.Errorf("status: got %d, want 200", result.StatusCode)
	}
}

func TestMemoryStore_Invalidate(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	_ = store.Store(ctx, "key2", 201, "created", 10*time.Minute)

	_ = store.Invalidate(ctx, "key2")

	result, _ := store.Check(ctx, "key2")
	if result != nil {
		t.Error("expected nil after invalidation")
	}
}

func TestCacheBackedStore_CheckMiss(t *testing.T) {
	c := newTestCache()
	store := NewStore(c)

	result, err := store.Check(context.Background(), "miss")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Error("expected nil for cache miss")
	}
}

func TestCacheBackedStore_StoreAndCheck(t *testing.T) {
	c := newTestCache()
	store := NewStore(c)
	ctx := context.Background()

	err := store.Store(ctx, "cached-key", 200, "response body", 10*time.Minute)
	if err != nil {
		t.Fatalf("store failed: %v", err)
	}

	result, err := store.Check(ctx, "cached-key")
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.StatusCode != 200 {
		t.Errorf("status: got %d, want 200", result.StatusCode)
	}
}

func TestCacheBackedStore_Invalidate(t *testing.T) {
	c := newTestCache()
	store := NewStore(c)
	ctx := context.Background()

	_ = store.Store(ctx, "inv-key", 200, "data", 10*time.Minute)
	_ = store.Invalidate(ctx, "inv-key")

	result, _ := store.Check(ctx, "inv-key")
	if result != nil {
		t.Error("expected nil after invalidation")
	}
}

// Helper: simple in-memory cache for tests

type testCache struct {
	mu    sync.RWMutex
	items map[string]testEntry
}

type testEntry struct {
	value     []byte
	expiresAt time.Time
}

func newTestCache() *testCache {
	return &testCache{items: make(map[string]testEntry)}
}

func (c *testCache) Get(_ context.Context, key string, dest interface{}) error {
	c.mu.RLock()
	entry, ok := c.items[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(entry.expiresAt) {
		return cache.ErrCacheMiss
	}
	return json.Unmarshal(entry.value, dest)
}

func (c *testCache) Set(_ context.Context, key string, value interface{}, ttl time.Duration) error {
	data, _ := json.Marshal(value)
	c.mu.Lock()
	c.items[key] = testEntry{value: data, expiresAt: time.Now().Add(ttl)}
	c.mu.Unlock()
	return nil
}

func (c *testCache) Delete(_ context.Context, key string) error {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
	return nil
}

func (c *testCache) Exists(_ context.Context, key string) (bool, error) {
	c.mu.RLock()
	_, ok := c.items[key]
	c.mu.RUnlock()
	return ok, nil
}

func (c *testCache) Close() error { return nil }
