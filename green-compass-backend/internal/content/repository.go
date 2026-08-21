package content

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("content not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Store(ctx context.Context, c *Content) error {
	query := `
		INSERT INTO place_content (
			place_id, generated_at, period_start, period_end,
			content_type, language, headline, body_text, call_to_action, source_indicators
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (place_id, content_type, period_start, period_end, language)
		DO UPDATE SET
			headline = $7,
			body_text = $8,
			call_to_action = $9,
			source_indicators = $10,
			generated_at = $2
		RETURNING id, created_at
	`

	err := r.pool.QueryRow(
		ctx, query,
		c.PlaceID, c.GeneratedAt, c.PeriodStart, c.PeriodEnd,
		c.ContentType, c.Language, c.Headline, c.BodyText, 		c.CallToAction, &c.SourceIndicators,
	).Scan(&c.ID, &c.CreatedAt)

	if err != nil {
		return fmt.Errorf("insert content: %w", err)
	}
	return nil
}

func (r *Repository) GetLatestByType(ctx context.Context, placeID uuid.UUID, contentType string) (*Content, error) {
	query := `
		SELECT id, place_id, generated_at, period_start, period_end,
		       content_type, language, headline, body_text, call_to_action, source_indicators, created_at
		FROM place_content
		WHERE place_id = $1 AND content_type = $2
		ORDER BY generated_at DESC
		LIMIT 1
	`
	return r.scanContent(ctx, query, placeID, contentType)
}

func (r *Repository) ListForPlace(ctx context.Context, placeID uuid.UUID, limit int) ([]Content, error) {
	query := `
		SELECT id, place_id, generated_at, period_start, period_end,
		       content_type, language, headline, body_text, call_to_action, source_indicators, created_at
		FROM place_content
		WHERE place_id = $1
		ORDER BY generated_at DESC
		LIMIT $2
	`
	rows, err := r.pool.Query(ctx, query, placeID, limit)
	if err != nil {
		return nil, fmt.Errorf("query content: %w", err)
	}
	defer rows.Close()

	var contents []Content
	for rows.Next() {
		var c Content
		if err := rows.Scan(&c.ID, &c.PlaceID, &c.GeneratedAt, &c.PeriodStart, &c.PeriodEnd,
			&c.ContentType, &c.Language, &c.Headline, &c.BodyText, &c.CallToAction, &c.SourceIndicators, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan content: %w", err)
		}
		contents = append(contents, c)
	}
	return contents, rows.Err()
}

func (r *Repository) scanContent(ctx context.Context, query string, args ...interface{}) (*Content, error) {
	var c Content

	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&c.ID, &c.PlaceID, &c.GeneratedAt, &c.PeriodStart, &c.PeriodEnd,
		&c.ContentType, &c.Language, &c.Headline, &c.BodyText, &c.CallToAction, &c.SourceIndicators, &c.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query content: %w", err)
	}

	return &c, nil
}
