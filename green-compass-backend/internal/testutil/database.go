package testutil

import (
	"context"
	"os"
	"testing"

	"green-compass-backend/internal/migrate"
	"green-compass-backend/migrations"
	"green-compass-backend/pkg/database"
)

func TestPool(t *testing.T) *database.Pool {
	t.Helper()
	url := os.Getenv("GC_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("GC_TEST_DATABASE_URL not set; skipping integration test")
	}
	pool, err := database.Connect(context.Background(), database.Options{URL: url})
	if err != nil {
		t.Fatalf("connect test database: %v", err)
	}
	t.Cleanup(pool.Close)

	if _, err := migrate.Up(context.Background(), pool, migrations.FS, 0); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
	return pool
}
