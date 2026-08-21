package httpx

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

const idempotencyHeader = "Idempotency-Key"

// keyStore is the minimal cache contract the idempotency middleware needs.
type keyStore interface {
	Check(c *gin.Context, key string) (bool, int, []byte, bool)
	Save(c *gin.Context, key string, status int, body []byte)
}

// memoryKeyStore is an in-process store for idempotent replays.
type memoryKeyStore struct {
	items map[string]replay
}

type replay struct {
	status int
	body   []byte
}

func newMemoryKeyStore() *memoryKeyStore {
	return &memoryKeyStore{items: make(map[string]replay)}
}

func (m *memoryKeyStore) Check(_ *gin.Context, key string) (bool, int, []byte, bool) {
	r, ok := m.items[key]
	return ok, r.status, r.body, ok
}

func (m *memoryKeyStore) Save(_ *gin.Context, key string, status int, body []byte) {
	m.items[key] = replay{status: status, body: body}
}

// IdempotencyMiddleware enforces Idempotency-Key semantics on mutating
// requests: a repeated key returns the stored response instead of
// re-executing the handler.
func IdempotencyMiddleware() gin.HandlerFunc {
	store := newMemoryKeyStore()
	return idempotencyWithStore(store)
}

func idempotencyWithStore(store keyStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		default:
			c.Next()
			return
		}

		key := c.GetHeader(idempotencyHeader)
		if key == "" {
			c.Next()
			return
		}
		cacheKey := c.Request.Method + ":" + c.FullPath() + ":" + key

		if hit, status, body, ok := store.Check(c, cacheKey); ok && hit {
			c.Data(status, "application/json", body)
			c.Abort()
			return
		}

		// Buffer the response so it can be replayed.
		writer := &bufferedWriter{ResponseWriter: c.Writer, buf: &bytes.Buffer{}}
		c.Writer = writer

		c.Next()

		if writer.status >= 200 && writer.status < 500 {
			store.Save(c, cacheKey, writer.status, writer.buf.Bytes())
		}
	}
}

type bufferedWriter struct {
	gin.ResponseWriter
	buf    *bytes.Buffer
	status int
}

func (w *bufferedWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *bufferedWriter) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	w.buf.Write(data)
	return w.ResponseWriter.Write(data)
}

var _ io.Writer = (*bufferedWriter)(nil)