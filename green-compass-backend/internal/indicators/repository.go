package indicators

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

var ErrNotFound = errors.New("indicator not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) GetDefinitionByCode(ctx context.Context, code string) (*Definition, error) {
	query := `
		SELECT id, code, display_name, category, source_variables, description, created_at, updated_at
		FROM indicator_definitions
		WHERE code = $1
	`
	var def Definition
	var variables pq.StringArray
	err := r.pool.QueryRow(ctx, query, code).Scan(
		&def.ID, &def.Code, &def.DisplayName, &def.Category, &variables, &def.Description, &def.CreatedAt, &def.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query definition: %w", err)
	}
	def.SourceVariables = variables
	return &def, nil
}

func (r *Repository) ListDefinitions(ctx context.Context) ([]Definition, error) {
	query := `
		SELECT id, code, display_name, category, source_variables, description, created_at, updated_at
		FROM indicator_definitions
		ORDER BY category, code
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query definitions: %w", err)
	}
	defer rows.Close()

	var defs []Definition
	for rows.Next() {
		var def Definition
		var variables pq.StringArray
		if err := rows.Scan(&def.ID, &def.Code, &def.DisplayName, &def.Category, &variables, &def.Description, &def.CreatedAt, &def.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan definition: %w", err)
		}
		def.SourceVariables = variables
		defs = append(defs, def)
	}
	return defs, rows.Err()
}

func (r *Repository) Store(ctx context.Context, ind *Indicator) error {
	metadataJSON, err := json.Marshal(ind.Metadata)
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}

	query := `
		INSERT INTO place_indicators (
			place_id, indicator_id, computed_at, period_start, period_end,
			value, unit, trend, trend_confidence, data_points_count, metadata
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (place_id, indicator_id, period_start, period_end)
		DO UPDATE SET
			value = $6,
			trend = $8,
			trend_confidence = $9,
			data_points_count = $10,
			metadata = $11,
			computed_at = $3
		RETURNING id, created_at
	`
	err = r.pool.QueryRow(
		ctx, query,
		ind.PlaceID, ind.IndicatorID, ind.ComputedAt, ind.PeriodStart, ind.PeriodEnd,
		ind.Value, ind.Unit, ind.Trend, ind.TrendConfidence, ind.DataPointsCount, metadataJSON,
	).Scan(&ind.ID, &ind.CreatedAt)

	if err != nil {
		return fmt.Errorf("insert indicator: %w", err)
	}
	return nil
}

func (r *Repository) GetLatest(ctx context.Context, placeID, indicatorID uuid.UUID) (*Indicator, error) {
	query := `
		SELECT id, place_id, indicator_id, computed_at, period_start, period_end,
		       value, unit, trend, trend_confidence, data_points_count, metadata, created_at
		FROM place_indicators
		WHERE place_id = $1 AND indicator_id = $2
		ORDER BY computed_at DESC
		LIMIT 1
	`
	return r.scanIndicator(ctx, query, placeID, indicatorID)
}

func (r *Repository) GetForPeriod(ctx context.Context, placeID, indicatorID uuid.UUID, periodStart, periodEnd time.Time) (*Indicator, error) {
	query := `
		SELECT id, place_id, indicator_id, computed_at, period_start, period_end,
		       value, unit, trend, trend_confidence, data_points_count, metadata, created_at
		FROM place_indicators
		WHERE place_id = $1 AND indicator_id = $2 AND period_start = $3 AND period_end = $4
		LIMIT 1
	`
	return r.scanIndicator(ctx, query, placeID, indicatorID, periodStart, periodEnd)
}

func (r *Repository) scanIndicator(ctx context.Context, query string, args ...interface{}) (*Indicator, error) {
	var ind Indicator
	var metadataJSON []byte
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&ind.ID, &ind.PlaceID, &ind.IndicatorID, &ind.ComputedAt, &ind.PeriodStart, &ind.PeriodEnd,
		&ind.Value, &ind.Unit, &ind.Trend, &ind.TrendConfidence, &ind.DataPointsCount, &metadataJSON, &ind.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query indicator: %w", err)
	}

	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &ind.Metadata); err != nil {
			return nil, fmt.Errorf("unmarshal metadata: %w", err)
		}
	} else {
		ind.Metadata = make(map[string]interface{})
	}
	return &ind, nil
}
