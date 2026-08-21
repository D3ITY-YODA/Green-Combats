package applicability

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrRuleNotFound = errors.New("applicability rule not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) GetRuleForIndicatorAndPlaceType(ctx context.Context, indicatorID uuid.UUID, placeType string) (*Rule, error) {
	query := `
		SELECT id, indicator_id, place_type, applicable, min_threshold, max_threshold, relevance_score, created_at, updated_at
		FROM indicator_applicability_rules
		WHERE indicator_id = $1 AND place_type = $2
	`
	var rule Rule
	err := r.pool.QueryRow(ctx, query, indicatorID, placeType).Scan(
		&rule.ID, &rule.IndicatorID, &rule.PlaceType, &rule.Applicable, &rule.MinThreshold, &rule.MaxThreshold, &rule.RelevanceScore, &rule.CreatedAt, &rule.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRuleNotFound
		}
		return nil, fmt.Errorf("query rule: %w", err)
	}
	return &rule, nil
}

func (r *Repository) ListRulesForPlaceType(ctx context.Context, placeType string) ([]Rule, error) {
	query := `
		SELECT id, indicator_id, place_type, applicable, min_threshold, max_threshold, relevance_score, created_at, updated_at
		FROM indicator_applicability_rules
		WHERE place_type = $1 AND applicable = TRUE
		ORDER BY relevance_score DESC
	`
	rows, err := r.pool.Query(ctx, query, placeType)
	if err != nil {
		return nil, fmt.Errorf("query rules: %w", err)
	}
	defer rows.Close()

	var rules []Rule
	for rows.Next() {
		var rule Rule
		if err := rows.Scan(&rule.ID, &rule.IndicatorID, &rule.PlaceType, &rule.Applicable, &rule.MinThreshold, &rule.MaxThreshold, &rule.RelevanceScore, &rule.CreatedAt, &rule.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan rule: %w", err)
		}
		rules = append(rules, rule)
	}
	return rules, rows.Err()
}
