package cache

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

var ErrCacheMiss = errors.New("cache miss")

// Cache provides key-value caching with optional TTL.
type Cache interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Close() error
}

// New returns a Cache backed by Redis when addr is configured,
// otherwise an in-memory cache for local development.
func New(addr string) Cache {
	if addr == "" {
		return NewMemoryCache()
	}
	return NewRedisCache(addr)
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

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{items: make(map[string]cacheEntry)}
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
	c.items[key] = cacheEntry{value: data, expiresAt: time.Now().Add(ttl)}
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

func (c *MemoryCache) Close() error { return nil }

// RedisCache is a minimal RESP-protocol Redis client supporting the
// command subset used by this service (GET / SET / DEL / EXISTS).
// It reconnects lazily and is safe for concurrent use.
type RedisCache struct {
	mu   sync.Mutex
	addr string
	conn net.Conn
	rw   *bufio.ReadWriter
}

func NewRedisCache(addr string) *RedisCache {
	return &RedisCache{addr: addr}
}

func (r *RedisCache) connect() error {
	if r.conn != nil {
		return nil
	}
	conn, err := net.DialTimeout("tcp", r.addr, 3*time.Second)
	if err != nil {
		return fmt.Errorf("redis dial %s: %w", r.addr, err)
	}
	r.conn = conn
	r.rw = bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	return nil
}

func (r *RedisCache) closeLocked() {
	if r.conn != nil {
		_ = r.conn.Close()
		r.conn = nil
		r.rw = nil
	}
}

// exec sends a command and returns the reply. One retry on connection failure.
func (r *RedisCache) exec(args ...string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	reply, err := r.execOnce(args)
	if err != nil && r.conn != nil {
		// Retry once with a fresh connection.
		r.closeLocked()
		reply, err = r.execOnce(args)
	}
	return reply, err
}

func (r *RedisCache) execOnce(args []string) (string, error) {
	if err := r.connect(); err != nil {
		return "", err
	}

	// Encode RESP array of bulk strings.
	cmd := "*" + strconv.Itoa(len(args)) + "\r\n"
	for _, a := range args {
		cmd += "$" + strconv.Itoa(len(a)) + "\r\n" + a + "\r\n"
	}
	if _, err := r.rw.WriteString(cmd); err != nil {
		return "", err
	}
	if err := r.rw.Flush(); err != nil {
		return "", err
	}

	return readReply(r.rw.Reader)
}

// readReply parses a single RESP reply into a string representation:
// simple strings and bulk strings return their payload, integers their
// decimal value, errors an error, and nil bulk "$-1" maps to ErrCacheMiss.
func readReply(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	line = line[:len(line)-2] // strip \r\n

	if len(line) == 0 {
		return "", fmt.Errorf("redis: empty reply")
	}

	switch line[0] {
	case '+': // simple string
		return line[1:], nil
	case '-': // error
		return "", fmt.Errorf("redis: %s", line[1:])
	case ':': // integer
		return line[1:], nil
	case '$': // bulk string
		n, err := strconv.Atoi(line[1:])
		if err != nil {
			return "", fmt.Errorf("redis: bad bulk length %q", line[1:])
		}
		if n == -1 {
			return "", ErrCacheMiss
		}
		buf := make([]byte, n+2) // payload + CRLF
		if _, err := ioReadFull(r, buf); err != nil {
			return "", err
		}
		return string(buf[:n]), nil
	default:
		return "", fmt.Errorf("redis: unknown reply type %q", line[0])
	}
}

func ioReadFull(r *bufio.Reader, buf []byte) (int, error) {
	total := 0
	for total < len(buf) {
		n, err := r.Read(buf[total:])
		total += n
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

func (r *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	reply, err := r.exec("GET", key)
	if err != nil {
		return err
	}
	if dest == nil {
		return nil
	}
	return json.Unmarshal([]byte(reply), dest)
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal value: %w", err)
	}
	_, err = r.exec("SET", key, string(data), "PX", strconv.FormatInt(ttl.Milliseconds(), 10))
	return err
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	_, err := r.exec("DEL", key)
	return err
}

func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	reply, err := r.exec("EXISTS", key)
	if err != nil {
		return false, err
	}
	return reply == "1", nil
}

func (r *RedisCache) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.closeLocked()
	return nil
}