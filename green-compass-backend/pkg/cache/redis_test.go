package cache

import (
	"context"
	"testing"
	"time"
)

func TestMemoryCache_SetAndGet(t *testing.T) {
	c := NewMemoryCache()
	ctx := context.Background()

	err := c.Set(ctx, "key1", "hello", 10*time.Second)
	if err != nil {
		t.Fatalf("set failed: %v", err)
	}

	var result string
	err = c.Get(ctx, "key1", &result)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if result != "hello" {
		t.Errorf("got %q, want %q", result, "hello")
	}
}

func TestMemoryCache_CacheMiss(t *testing.T) {
	c := NewMemoryCache()
	ctx := context.Background()

	var result string
	err := c.Get(ctx, "nonexistent", &result)
	if err != ErrCacheMiss {
		t.Errorf("expected ErrCacheMiss, got %v", err)
	}
}

func TestMemoryCache_TTLExpiration(t *testing.T) {
	c := NewMemoryCache()
	ctx := context.Background()

	err := c.Set(ctx, "expire", "value", 50*time.Millisecond)
	if err != nil {
		t.Fatalf("set failed: %v", err)
	}

	// Should exist immediately
	var result string
	err = c.Get(ctx, "expire", &result)
	if err != nil {
		t.Fatalf("get should succeed immediately: %v", err)
	}

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	err = c.Get(ctx, "expire", &result)
	if err != ErrCacheMiss {
		t.Errorf("expected ErrCacheMiss after expiration, got %v", err)
	}
}

func TestMemoryCache_Delete(t *testing.T) {
	c := NewMemoryCache()
	ctx := context.Background()

	_ = c.Set(ctx, "del", "value", 10*time.Second)

	exists, _ := c.Exists(ctx, "del")
	if !exists {
		t.Fatal("key should exist before delete")
	}

	_ = c.Delete(ctx, "del")

	exists, _ = c.Exists(ctx, "del")
	if exists {
		t.Fatal("key should not exist after delete")
	}
}

func TestMemoryCache_Exists(t *testing.T) {
	c := NewMemoryCache()
	ctx := context.Background()

	exists, _ := c.Exists(ctx, "exists?")
	if exists {
		t.Fatal("should not exist before set")
	}

	_ = c.Set(ctx, "exists?", true, 10*time.Second)

	exists, _ = c.Exists(ctx, "exists?")
	if !exists {
		t.Fatal("should exist after set")
	}
}

func TestMemoryCache_Overwrite(t *testing.T) {
	c := NewMemoryCache()
	ctx := context.Background()

	_ = c.Set(ctx, "key", "first", 10*time.Second)
	_ = c.Set(ctx, "key", "second", 10*time.Second)

	var result string
	_ = c.Get(ctx, "key", &result)
	if result != "second" {
		t.Errorf("got %q, want %q", result, "second")
	}
}

func TestMemoryCache_Struct(t *testing.T) {
	c := NewMemoryCache()
	ctx := context.Background()

	type item struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	_ = c.Set(ctx, "struct", item{Name: "test", Value: 42}, 10*time.Second)

	var got item
	_ = c.Get(ctx, "struct", &got)
	if got.Name != "test" || got.Value != 42 {
		t.Errorf("got %+v, want {test 42}", got)
	}
}

func TestRedisCache_FallbackToMemory(t *testing.T) {
	c := NewRedisCache("") // empty addr falls back to memory
	ctx := context.Background()

	_ = c.Set(ctx, "fallback", "value", 10*time.Second)

	var result string
	_ = c.Get(ctx, "fallback", &result)
	if result != "value" {
		t.Errorf("got %q, want %q", result, "value")
	}
}

func TestMemoryCache_Close(t *testing.T) {
	c := NewMemoryCache()
	if err := c.Close(); err != nil {
		t.Errorf("close failed: %v", err)
	}
}
