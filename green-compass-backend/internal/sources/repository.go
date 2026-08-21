package sources

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"green-compass-backend/pkg/database"
)

type Repository struct{ pool *database.Pool }

func NewRepository(pool *database.Pool) *Repository { return &Repository{pool: pool} }

func (r *Repository) ByCode(ctx context.Context, code string) (*Source, error) {
	source := &Source{}
	err := scanSource(r.pool.QueryRow(ctx, `
		SELECT id, code, display_name, poll_interval_seconds, enabled, config, created_at, updated_at
		FROM data_sources WHERE code = $1`, code), source)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return source, nil
}

func (r *Repository) Enabled(ctx context.Context) ([]Source, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, code, display_name, poll_interval_seconds, enabled, config, created_at, updated_at
		FROM data_sources WHERE enabled ORDER BY code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Source
	for rows.Next() {
		var source Source
		var seconds int64
		if err := rows.Scan(&source.ID, &source.Code, &source.DisplayName, &seconds, &source.Enabled, &source.Config, &source.CreatedAt, &source.UpdatedAt); err != nil {
			return nil, err
		}
		source.PollInterval = time.Duration(seconds) * time.Second
		result = append(result, source)
	}
	return result, rows.Err()
}

type sourceScanner interface{ Scan(...any) error }

func scanSource(row sourceScanner, source *Source) error {
	var seconds int64
	if err := row.Scan(&source.ID, &source.Code, &source.DisplayName, &seconds, &source.Enabled, &source.Config, &source.CreatedAt, &source.UpdatedAt); err != nil {
		return err
	}
	source.PollInterval = time.Duration(seconds) * time.Second
	return nil
}
