package assessment

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("assessment not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Store(ctx context.Context, a *Assessment) error {
	query := `
		INSERT INTO place_assessments (
			place_id, assessed_at, period_start, period_end,
			urgency_score, confidence_score, applicable_indicators, affected_groups, assessment_summary
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (place_id, period_start, period_end)
		DO UPDATE SET
			urgency_score = $5,
			confidence_score = $6,
			applicable_indicators = $7,
			affected_groups = $8,
			assessment_summary = $9,
			assessed_at = $2
		RETURNING id, created_at
	`

	err := r.pool.QueryRow(
		ctx, query,
		a.PlaceID, a.AssessedAt, a.PeriodStart, a.PeriodEnd,
		a.UrgencyScore, a.ConfidenceScore, a.ApplicableIndicators, a.AffectedGroups, a.AssessmentSummary,
	).Scan(&a.ID, &a.CreatedAt)

	if err != nil {
		return fmt.Errorf("insert assessment: %w", err)
	}
	return nil
}

func (r *Repository) GetLatest(ctx context.Context, placeID uuid.UUID) (*Assessment, error) {
	query := `
		SELECT id, place_id, assessed_at, period_start, period_end,
		       urgency_score, confidence_score, applicable_indicators, affected_groups, assessment_summary, created_at
		FROM place_assessments
		WHERE place_id = $1
		ORDER BY assessed_at DESC
		LIMIT 1
	`
	return r.scanAssessment(ctx, query, placeID)
}

func (r *Repository) GetForPeriod(ctx context.Context, placeID uuid.UUID, periodStart, periodEnd time.Time) (*Assessment, error) {
	query := `
		SELECT id, place_id, assessed_at, period_start, period_end,
		       urgency_score, confidence_score, applicable_indicators, affected_groups, assessment_summary, created_at
		FROM place_assessments
		WHERE place_id = $1 AND period_start = $2 AND period_end = $3
		LIMIT 1
	`
	var a Assessment

	err := r.pool.QueryRow(ctx, query, placeID, periodStart, periodEnd).Scan(
		&a.ID, &a.PlaceID, &a.AssessedAt, &a.PeriodStart, &a.PeriodEnd,
		&a.UrgencyScore, &a.ConfidenceScore, &a.ApplicableIndicators, &a.AffectedGroups, &a.AssessmentSummary, &a.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query assessment: %w", err)
	}

	return &a, nil
}

func (r *Repository) scanAssessment(ctx context.Context, query string, args ...interface{}) (*Assessment, error) {
	var a Assessment

	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&a.ID, &a.PlaceID, &a.AssessedAt, &a.PeriodStart, &a.PeriodEnd,
		&a.UrgencyScore, &a.ConfidenceScore, &a.ApplicableIndicators, &a.AffectedGroups, &a.AssessmentSummary, &a.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query assessment: %w", err)
	}

	return &a, nil
}

