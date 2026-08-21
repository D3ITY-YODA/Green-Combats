package migrate

import (
	"cmp"
	"context"
	"fmt"
	"io/fs"
	"regexp"
	"slices"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"

	"green-compass-backend/pkg/database"
)

const advisoryLockKey int64 = 724918233401

var fileNameRe = regexp.MustCompile(`^([0-9]+)_([a-z0-9_]+)\.(up|down)\.sql$`)

type migrationFile struct {
	Version int64
	Name    string
	Dir     string
	Path    string
}

type Status struct {
	Current int64
	Pending []int64
}

func parseMigrations(fsys fs.FS) ([]migrationFile, error) {
	entries, err := fs.Glob(fsys, "*.sql")
	if err != nil {
		return nil, fmt.Errorf("scan migrations: %w", err)
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("scan migrations: no .sql files found")
	}

	var files []migrationFile
	for _, entry := range entries {
		m := fileNameRe.FindStringSubmatch(entry)
		if m == nil {
			return nil, fmt.Errorf("invalid migration filename %q: want NNN_name.up.sql or NNN_name.down.sql", entry)
		}
		version, err := strconv.ParseInt(m[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid migration version in %q: %w", entry, err)
		}
		files = append(files, migrationFile{Version: version, Name: m[2], Dir: m[3], Path: entry})
	}

	slices.SortFunc(files, func(a, b migrationFile) int {
		if c := cmp.Compare(a.Version, b.Version); c != 0 {
			return c
		}
		return cmp.Compare(a.Dir, b.Dir)
	})

	type pair struct{ up, down *migrationFile }
	byVersion := map[int64]*pair{}
	for i := range files {
		f := &files[i]
		p, ok := byVersion[f.Version]
		if !ok {
			p = &pair{}
			byVersion[f.Version] = p
		}
		switch f.Dir {
		case "up":
			if p.up != nil {
				return nil, fmt.Errorf("duplicate up migration for version %d (%q and %q)", f.Version, p.up.Name, f.Name)
			}
			p.up = f
		case "down":
			if p.down != nil {
				return nil, fmt.Errorf("duplicate down migration for version %d (%q and %q)", f.Version, p.down.Name, f.Name)
			}
			p.down = f
		}
	}

	for version, p := range byVersion {
		if p.up == nil || p.down == nil {
			name := "<unknown>"
			if p.up != nil {
				name = p.up.Name
			} else if p.down != nil {
				name = p.down.Name
			}
			return nil, fmt.Errorf("version %d (%s) must have both .up.sql and .down.sql", version, name)
		}
		if p.up.Name != p.down.Name {
			return nil, fmt.Errorf("version %d name mismatch: up=%q down=%q", version, p.up.Name, p.down.Name)
		}
	}

	return files, nil
}

func Up(ctx context.Context, pool *database.Pool, fsys fs.FS, steps int) ([]int64, error) {
	conn, unlock, err := acquire(ctx, pool)
	if err != nil {
		return nil, err
	}
	defer unlock()

	files, err := parseMigrations(fsys)
	if err != nil {
		return nil, err
	}
	current, err := currentVersion(ctx, conn)
	if err != nil {
		return nil, err
	}

	var applied []int64
	for _, f := range files {
		if f.Dir != "up" || f.Version <= current {
			continue
		}
		if steps > 0 && len(applied) >= steps {
			break
		}
		if err := applyOne(ctx, conn, fsys, f, true); err != nil {
			return applied, fmt.Errorf("apply migration %06d_%s: %w", f.Version, f.Name, err)
		}
		applied = append(applied, f.Version)
	}
	return applied, nil
}

func Down(ctx context.Context, pool *database.Pool, fsys fs.FS, steps int) ([]int64, error) {
	conn, unlock, err := acquire(ctx, pool)
	if err != nil {
		return nil, err
	}
	defer unlock()

	files, err := parseMigrations(fsys)
	if err != nil {
		return nil, err
	}
	current, err := currentVersion(ctx, conn)
	if err != nil {
		return nil, err
	}

	var reverted []int64
	for i := len(files) - 1; i >= 0; i-- {
		f := files[i]
		if f.Dir != "down" || f.Version > current {
			continue
		}
		if steps > 0 && len(reverted) >= steps {
			break
		}
		if err := applyOne(ctx, conn, fsys, f, false); err != nil {
			return reverted, fmt.Errorf("revert migration %06d_%s: %w", f.Version, f.Name, err)
		}
		reverted = append(reverted, f.Version)
	}
	return reverted, nil
}

func Current(ctx context.Context, pool *database.Pool, fsys fs.FS) (Status, error) {
	conn, unlock, err := acquire(ctx, pool)
	if err != nil {
		return Status{}, err
	}
	defer unlock()

	files, err := parseMigrations(fsys)
	if err != nil {
		return Status{}, err
	}
	current, err := currentVersion(ctx, conn)
	if err != nil {
		return Status{}, err
	}

	st := Status{Current: current}
	for _, f := range files {
		if f.Dir == "up" && f.Version > current {
			st.Pending = append(st.Pending, f.Version)
		}
	}
	return st, nil
}

func acquire(ctx context.Context, pool *database.Pool) (*pgxpool.Conn, func(), error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("acquire connection: %w", err)
	}
	if _, err := conn.Exec(ctx, "SELECT pg_advisory_lock($1)", advisoryLockKey); err != nil {
		conn.Release()
		return nil, nil, fmt.Errorf("acquire migration lock: %w", err)
	}
	unlock := func() {
		_, _ = conn.Exec(context.WithoutCancel(ctx), "SELECT pg_advisory_unlock($1)", advisoryLockKey)
		conn.Release()
	}
	return conn, unlock, nil
}

func currentVersion(ctx context.Context, conn *pgxpool.Conn) (int64, error) {
	if _, err := conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
		    version    BIGINT PRIMARY KEY,
		    name       TEXT NOT NULL,
		    applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)`); err != nil {
		return 0, fmt.Errorf("ensure schema_migrations table: %w", err)
	}

	var current int64
	if err := conn.QueryRow(ctx, "SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&current); err != nil {
		return 0, fmt.Errorf("read current migration version: %w", err)
	}
	return current, nil
}

func applyOne(ctx context.Context, conn *pgxpool.Conn, fsys fs.FS, f migrationFile, up bool) error {
	sqlBytes, err := fs.ReadFile(fsys, f.Path)
	if err != nil {
		return fmt.Errorf("read migration file: %w", err)
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if _, err := tx.Exec(ctx, string(sqlBytes)); err != nil {
		return err
	}
	if up {
		if _, err := tx.Exec(ctx,
			"INSERT INTO schema_migrations (version, name) VALUES ($1, $2)", f.Version, f.Name); err != nil {
			return err
		}
	} else {
		if _, err := tx.Exec(ctx,
			"DELETE FROM schema_migrations WHERE version = $1", f.Version); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
