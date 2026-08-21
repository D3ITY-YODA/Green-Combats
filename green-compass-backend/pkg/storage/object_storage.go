package storage

import (
	"context"
	"io"
	"sync"
)

type Store interface {
	Put(ctx context.Context, key string, r io.Reader, contentType string) (string, error)
	Delete(ctx context.Context, key string) error
}

type memoryStore struct {
	mu    sync.RWMutex
	data  map[string][]byte
	types map[string]string
}

func NewMemoryStore() Store {
	return &memoryStore{
		data:  make(map[string][]byte),
		types: make(map[string]string),
	}
}

func (m *memoryStore) Put(_ context.Context, key string, r io.Reader, contentType string) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	m.mu.Lock()
	m.data[key] = data
	m.types[key] = contentType
	m.mu.Unlock()
	return key, nil
}

func (m *memoryStore) Delete(_ context.Context, key string) error {
	m.mu.Lock()
	delete(m.data, key)
	delete(m.types, key)
	m.mu.Unlock()
	return nil
}
