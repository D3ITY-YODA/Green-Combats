package updates

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("update not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// GetTodayForPlace returns the latest "today" content for a place
func (r *Repository) GetTodayForPlace(ctx context.Context, placeID uuid.UUID) (*Update, error) {
	query := `
		SELECT c.id, c.place_id, p.name, c.content_type, c.headline, c.body_text, c.call_to_action,
		       c.generated_at, c.period_start, c.period_end
		FROM place_content c
		JOIN places p ON c.place_id = p.id
		WHERE c.place_id = $1 AND c.content_type = 'today'
		ORDER BY c.generated_at DESC
		LIMIT 1
	`

	var update Update
	var callToAction *string

	err := r.pool.QueryRow(ctx, query, placeID).Scan(
		&update.ID, &update.PlaceID, &update.PlaceName, &update.ContentType, &update.Headline, &update.BodyText, &callToAction,
		&update.GeneratedAt, &update.PeriodStart, &update.PeriodEnd,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query today: %w", err)
	}

	update.CallToAction = callToAction
	return &update, nil
}

// ListForUser returns all updates for user's saved places, paginated
func (r *Repository) ListForUser(ctx context.Context, req ListRequest) ([]Update, int, error) {
	// Build query to get content from places user has saved
	baseQuery := `
		SELECT c.id, c.place_id, p.name, c.content_type, c.headline, c.body_text, c.call_to_action,
		       c.generated_at, c.period_start, c.period_end
		FROM place_content c
		JOIN places p ON c.place_id = p.id
		JOIN user_saved_places usp ON usp.place_id = p.id
		WHERE usp.user_id = $1
	`

	args := []interface{}{req.UserID}
	argIdx := 2

	// Filter by specific place if provided
	if req.PlaceID != nil {
		baseQuery += fmt.Sprintf(" AND c.place_id = $%d", argIdx)
		args = append(args, *req.PlaceID)
		argIdx++
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(DISTINCT c.id) FROM (%s) AS subq", baseQuery)
	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count updates: %w", err)
	}

	// Get paginated results, latest first
	offset := (req.Page - 1) * req.Limit
	listQuery := baseQuery + fmt.Sprintf(`
		ORDER BY c.generated_at DESC
		LIMIT $%d OFFSET $%d
	`, argIdx, argIdx+1)

	args = append(args, req.Limit, offset)

	rows, err := r.pool.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query updates: %w", err)
	}
	defer rows.Close()

	var updates []Update
	for rows.Next() {
		var update Update
		var callToAction *string

		if err := rows.Scan(&update.ID, &update.PlaceID, &update.PlaceName, &update.ContentType, &update.Headline, &update.BodyText, &callToAction,
			&update.GeneratedAt, &update.PeriodStart, &update.PeriodEnd); err != nil {
			return nil, 0, fmt.Errorf("scan update: %w", err)
		}
		update.CallToAction = callToAction
		updates = append(updates, update)
	}

	return updates, total, rows.Err()
}

// ListIndicators returns indicators for a given content
func (r *Repository) ListIndicators(ctx context.Context, contentID uuid.UUID) ([]Indicator, error) {
	query := `
		SELECT ind_def.id, ind_def.code, ind_def.display_name, 
		       place_ind.value, place_ind.unit, place_ind.trend,
		       iar.relevance_score
		FROM place_indicators place_ind
		JOIN indicator_definitions ind_def ON place_ind.indicator_id = ind_def.id
		JOIN place_content pc ON pc.place_id = place_ind.place_id
		LEFT JOIN indicator_applicability_rules iar 
			ON iar.indicator_id = ind_def.id AND iar.applicable = TRUE
		WHERE pc.id = $1
		ORDER BY iar.relevance_score DESC NULLS LAST
	`

	rows, err := r.pool.Query(ctx, query, contentID)
	if err != nil {
		return nil, fmt.Errorf("query indicators: %w", err)
	}
	defer rows.Close()

	var indicators []Indicator
	for rows.Next() {
		var ind Indicator
		var relScore *float64

		if err := rows.Scan(&ind.ID, &ind.Code, &ind.DisplayName, &ind.Value, &ind.Unit, &ind.Trend, &relScore); err != nil {
			return nil, fmt.Errorf("scan indicator: %w", err)
		}
		if relScore != nil {
			ind.RelevanceScore = *relScore
		} else {
			ind.RelevanceScore = 1.0
		}
		indicators = append(indicators, ind)
	}

	return indicators, rows.Err()
}

// ExploreIndicators returns computed indicators for exploration, with optional filters
func (r *Repository) ExploreIndicators(ctx context.Context, req ExploreRequest) ([]ExploreIndicator, int, error) {
	baseQuery := `
		SELECT pi.id, pi.place_id, p.name, ind_def.code, ind_def.display_name, ind_def.category,
		       pi.value, pi.unit, pi.trend, pi.trend_confidence, pi.computed_at
		FROM place_indicators pi
		JOIN places p ON pi.place_id = p.id
		JOIN indicator_definitions ind_def ON pi.indicator_id = ind_def.id
		WHERE pi.computed_at > NOW() - INTERVAL '7 days'
	`

	args := []interface{}{}
	argIdx := 1

	// Filter by category if provided
	if req.Category != nil {
		baseQuery += fmt.Sprintf(" AND ind_def.category = $%d", argIdx)
		args = append(args, *req.Category)
		argIdx++
	}

	// Filter by specific place if provided
	if req.PlaceID != nil {
		baseQuery += fmt.Sprintf(" AND pi.place_id = $%d", argIdx)
		args = append(args, *req.PlaceID)
		argIdx++
	}

	// Otherwise, for discovery: show only user's saved places or nearby places
	// For now, show all; this can be scoped per user saved places in production
	baseQuery += fmt.Sprintf(`
		ORDER BY pi.computed_at DESC
	`)

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(DISTINCT pi.id) FROM (%s) AS subq", baseQuery)
	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count indicators: %w", err)
	}

	// Get paginated results
	offset := (req.Page - 1) * req.Limit
	listQuery := baseQuery + fmt.Sprintf(`
		LIMIT $%d OFFSET $%d
	`, argIdx+1, argIdx+2)

	args = append(args, req.Limit, offset)

	rows, err := r.pool.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query explore: %w", err)
	}
	defer rows.Close()

	var indicators []ExploreIndicator
	for rows.Next() {
		var ind ExploreIndicator
		var trend *string
		var trendConfidence *float64

		if err := rows.Scan(&ind.ID, &ind.PlaceID, &ind.PlaceName, &ind.Code, &ind.DisplayName, &ind.Category,
			&ind.Value, &ind.Unit, &trend, &trendConfidence, &ind.ComputedAt); err != nil {
			return nil, 0, fmt.Errorf("scan explore indicator: %w", err)
		}
		ind.Trend = trend
		ind.TrendConfidence = trendConfidence
		indicators = append(indicators, ind)
	}

	return indicators, total, rows.Err()
}
