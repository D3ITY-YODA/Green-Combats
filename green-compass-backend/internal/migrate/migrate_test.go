package migrate

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"green-compass-backend/pkg/database"
)

func isolatedPool(t *testing.T) *database.Pool {
	t.Helper()
	base := os.Getenv("GC_TEST_DATABASE_URL")
	if base == "" {
		t.Skip("GC_TEST_DATABASE_URL not set; skipping integration test")
	}
	ctx := context.Background()

	admin, err := database.Connect(ctx, database.Options{URL: base})
	if err != nil {
		t.Fatalf("connect admin database: %v", err)
	}
	t.Cleanup(admin.Close)

	name := fmt.Sprintf("gc_migrate_%d", time.Now().UnixNano())
	if _, err := admin.Exec(ctx, "CREATE DATABASE "+name); err != nil {
		t.Fatalf("create isolated database: %v", err)
	}
	t.Cleanup(func() {
		_, _ = admin.Exec(ctx, "DROP DATABASE "+name+" WITH (FORCE)")
	})

	parsed, err := url.Parse(base)
	if err != nil {
		t.Fatalf("parse test database url: %v", err)
	}
	parsed.Path = "/" + name

	pool, err := database.Connect(ctx, database.Options{URL: parsed.String()})
	if err != nil {
		t.Fatalf("connect isolated database: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

func validFS() fstest.MapFS {
	return fstest.MapFS{
		"001_a.up.sql":   {Data: []byte("SELECT 1;")},
		"001_a.down.sql": {Data: []byte("SELECT 1;")},
		"002_b.up.sql":   {Data: []byte("SELECT 2;")},
		"002_b.down.sql": {Data: []byte("SELECT 2;")},
		"010_c.up.sql":   {Data: []byte("SELECT 3;")},
		"010_c.down.sql": {Data: []byte("SELECT 3;")},
	}
}

func TestParseMigrations_Valid(t *testing.T) {
	files, err := parseMigrations(validFS())
	if err != nil {
		t.Fatalf("parseMigrations() unexpected error: %v", err)
	}

	wantVersions := []int64{1, 1, 2, 2, 10, 10}
	for i, f := range files {
		if f.Version != wantVersions[i] {
			t.Errorf("files[%d].Version = %d, want %d", i, f.Version, wantVersions[i])
		}
	}
	if files[0].Dir != "down" || files[1].Dir != "up" {
		t.Errorf("within a version, down should sort before up: got %q then %q", files[0].Dir, files[1].Dir)
	}
}

func TestParseMigrations_Errors(t *testing.T) {
	tests := []struct {
		name    string
		fsys    fstest.MapFS
		wantErr string
	}{
		{
			name: "missing down file",
			fsys: fstest.MapFS{
				"001_a.up.sql": {Data: []byte("SELECT 1;")},
			},
			wantErr: "must have both",
		},
		{
			name: "missing up file",
			fsys: fstest.MapFS{
				"001_a.down.sql": {Data: []byte("SELECT 1;")},
			},
			wantErr: "must have both",
		},
		{
			name: "invalid filename",
			fsys: fstest.MapFS{
				"001-a.up.sql":   {Data: []byte("SELECT 1;")},
				"001-a.down.sql": {Data: []byte("SELECT 1;")},
			},
			wantErr: "invalid migration filename",
		},
		{
			name: "duplicate version different names",
			fsys: fstest.MapFS{
				"001_a.up.sql":   {Data: []byte("SELECT 1;")},
				"001_a.down.sql": {Data: []byte("SELECT 1;")},
				"001_b.up.sql":   {Data: []byte("SELECT 2;")},
				"001_b.down.sql": {Data: []byte("SELECT 2;")},
			},
			wantErr: "duplicate down migration",
		},
		{
			name:    "empty filesystem",
			fsys:    fstest.MapFS{},
			wantErr: "no .sql files",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseMigrations(tt.fsys)
			if err == nil {
				t.Fatal("parseMigrations() expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error = %v, want it to contain %q", err, tt.wantErr)
			}
		})
	}
}

func TestUpDownCycle_Integration(t *testing.T) {
	pool := isolatedPool(t)
	ctx := context.Background()
	fsys := os.DirFS("../../migrations")

	if _, err := Down(ctx, pool, fsys, 0); err != nil {
		t.Fatalf("reset to clean slate: %v", err)
	}

	st, err := Current(ctx, pool, fsys)
	if err != nil {
		t.Fatalf("Status() unexpected error: %v", err)
	}
	if st.Current != 0 {
		t.Fatalf("Current = %d after reset, want 0", st.Current)
	}
	if len(st.Pending) != 10 {
		t.Fatalf("Pending = %v, want 10 versions", st.Pending)
	}

	applied, err := Up(ctx, pool, fsys, 0)
	if err != nil {
		t.Fatalf("Up(all) unexpected error: %v", err)
	}
	if len(applied) != 10 {
		t.Fatalf("applied = %v, want 10 versions", applied)
	}

	st, err = Current(ctx, pool, fsys)
	if err != nil {
		t.Fatalf("Status() unexpected error: %v", err)
	}
	if st.Current != 10 || len(st.Pending) != 0 {
		t.Fatalf("status after up-all = current %d pending %v, want current 10 pending empty", st.Current, st.Pending)
	}

	for _, table := range []string{"users", "organizations", "organization_members", "places", "user_saved_places", "observations", "updates", "data_sources", "ingestion_runs", "raw_records"} {
		var exists bool
		if err := pool.QueryRow(ctx,
			"SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)", table).Scan(&exists); err != nil {
			t.Fatalf("check table %s: %v", table, err)
		}
		if !exists {
			t.Errorf("table %s does not exist after up-all", table)
		}
	}

	reverted, err := Down(ctx, pool, fsys, 2)
	if err != nil {
		t.Fatalf("Down(2) unexpected error: %v", err)
	}
	if len(reverted) != 2 {
		t.Fatalf("reverted = %v, want 2 versions", reverted)
	}
	st, _ = Current(ctx, pool, fsys)
	if st.Current != 8 {
		t.Fatalf("Current = %d after Down(2), want 8", st.Current)
	}

	if _, err := Up(ctx, pool, fsys, 0); err != nil {
		t.Fatalf("re-Up(all) unexpected error: %v", err)
	}
	st, _ = Current(ctx, pool, fsys)
	if st.Current != 10 {
		t.Fatalf("Current = %d after re-up, want 10", st.Current)
	}

	if _, err := Up(ctx, pool, fsys, 0); err != nil {
		t.Fatalf("idempotent Up(all) unexpected error: %v", err)
	}

	if _, err := Down(ctx, pool, fsys, 0); err != nil {
		t.Fatalf("Down(all) unexpected error: %v", err)
	}
	st, _ = Current(ctx, pool, fsys)
	if st.Current != 0 {
		t.Fatalf("Current = %d after down-all, want 0", st.Current)
	}

	for _, table := range []string{"users", "places", "updates"} {
		var exists bool
		if err := pool.QueryRow(ctx,
			"SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)", table).Scan(&exists); err != nil {
			t.Fatalf("check table %s: %v", table, err)
		}
		if exists {
			t.Errorf("table %s still exists after down-all", table)
		}
	}
}

func TestDown_AtZeroIsNoOp_Integration(t *testing.T) {
	pool := isolatedPool(t)
	ctx := context.Background()
	fsys := os.DirFS("../../migrations")

	if _, err := Down(ctx, pool, fsys, 0); err != nil {
		t.Fatalf("Down at zero: %v", err)
	}

	reverted, err := Down(ctx, pool, fsys, 0)
	if err != nil {
		t.Fatalf("Down(all) at zero unexpected error: %v", err)
	}
	if len(reverted) != 0 {
		t.Fatalf("reverted = %v, want none at zero", reverted)
	}
}
