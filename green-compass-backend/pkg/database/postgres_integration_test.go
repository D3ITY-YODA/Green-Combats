package database_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"green-compass-backend/pkg/database"
)

var scratchMu sync.Mutex

func testPool(t *testing.T) *database.Pool {
	t.Helper()
	url := os.Getenv("GC_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("GC_TEST_DATABASE_URL not set; skipping integration test")
	}
	pool, err := database.Connect(context.Background(), database.Options{URL: url})
	if err != nil {
		t.Fatalf("Connect() unexpected error: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

func scratchTable(t *testing.T, pool *database.Pool) string {
	t.Helper()
	scratchMu.Lock()
	defer scratchMu.Unlock()

	name := fmt.Sprintf("dbtest_%d", time.Now().UnixNano())
	_, err := pool.Exec(context.Background(), fmt.Sprintf(
		`CREATE TABLE %s (id INT PRIMARY KEY, note TEXT)`, name))
	if err != nil {
		t.Fatalf("create scratch table: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), fmt.Sprintf(`DROP TABLE IF EXISTS %s`, name))
	})
	return name
}

func rowCount(t *testing.T, pool *database.Pool, table string) int {
	t.Helper()
	var n int
	err := pool.QueryRow(context.Background(),
		fmt.Sprintf(`SELECT COUNT(*) FROM %s`, table)).Scan(&n)
	if err != nil {
		t.Fatalf("count rows: %v", err)
	}
	return n
}

func insertRow(ctx context.Context, tx pgx.Tx, table string, id int) error {
	_, err := tx.Exec(ctx, fmt.Sprintf(`INSERT INTO %s (id, note) VALUES ($1, $2)`, table), id, "x")
	return err
}

func TestConnect_Integration(t *testing.T) {
	pool := testPool(t)

	if err := pool.Ping(context.Background()); err != nil {
		t.Fatalf("Ping() unexpected error: %v", err)
	}
}

func TestConnect_Errors(t *testing.T) {
	tests := []struct {
		name string
		opts database.Options
	}{
		{name: "empty url", opts: database.Options{}},
		{name: "unparseable url", opts: database.Options{URL: "://not-a-url"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := database.Connect(context.Background(), tt.opts); err == nil {
				t.Fatalf("Connect(%+v) expected error, got nil", tt.opts)
			}
		})
	}
}

func TestWithTx_CommitsOnSuccess(t *testing.T) {
	pool := testPool(t)
	table := scratchTable(t, pool)

	err := database.WithTx(context.Background(), pool, func(ctx context.Context, tx pgx.Tx) error {
		return insertRow(ctx, tx, table, 1)
	})
	if err != nil {
		t.Fatalf("WithTx() unexpected error: %v", err)
	}

	if got := rowCount(t, pool, table); got != 1 {
		t.Fatalf("row count = %d, want 1 after committed tx", got)
	}
}

func TestWithTx_RollsBackOnError(t *testing.T) {
	pool := testPool(t)
	table := scratchTable(t, pool)

	wantErr := fmt.Errorf("boom")
	err := database.WithTx(context.Background(), pool, func(ctx context.Context, tx pgx.Tx) error {
		if err := insertRow(ctx, tx, table, 1); err != nil {
			return err
		}
		return wantErr
	})
	if err == nil || err.Error() != wantErr.Error() {
		t.Fatalf("WithTx() error = %v, want %v", err, wantErr)
	}

	if got := rowCount(t, pool, table); got != 0 {
		t.Fatalf("row count = %d, want 0 after rolled back tx", got)
	}
}

func TestWithTx_RollsBackOnPanic(t *testing.T) {
	pool := testPool(t)
	table := scratchTable(t, pool)

	func() {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("expected panic to propagate out of WithTx")
			}
		}()
		_ = database.WithTx(context.Background(), pool, func(ctx context.Context, tx pgx.Tx) error {
			if err := insertRow(ctx, tx, table, 1); err != nil {
				return err
			}
			panic("kaboom")
		})
	}()

	if got := rowCount(t, pool, table); got != 0 {
		t.Fatalf("row count = %d, want 0 after panicked tx", got)
	}
}

func TestWithTx_QueryInsideTxSeesOwnWrites(t *testing.T) {
	pool := testPool(t)
	table := scratchTable(t, pool)

	err := database.WithTx(context.Background(), pool, func(ctx context.Context, tx pgx.Tx) error {
		if err := insertRow(ctx, tx, table, 7); err != nil {
			return err
		}
		var n int
		if err := tx.QueryRow(ctx, fmt.Sprintf(`SELECT COUNT(*) FROM %s`, table)).Scan(&n); err != nil {
			return err
		}
		if n != 1 {
			t.Errorf("rows visible inside tx = %d, want 1", n)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("WithTx() unexpected error: %v", err)
	}
}
