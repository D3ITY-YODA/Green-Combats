package reporting

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// DeliveryStats queries notification delivery metrics for a place within a time range.
func (r *Repository) DeliveryStats(ctx context.Context, placeID uuid.UUID, from, to string) (*DeliveryStats, error) {
	stats := &DeliveryStats{
		PlaceID:     placeID,
		ByChannel:   make(map[string]int),
		PeriodStart: parseTime(from),
		PeriodEnd:   parseTime(to),
	}

	// Count total sent
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM notifications n
		JOIN users u ON n.user_id = u.id
		JOIN user_saved_places usp ON usp.user_id = u.id
		WHERE usp.place_id = $1 AND n.status = 'sent'
		AND n.created_at BETWEEN $2 AND $3
	`, placeID, from, to).Scan(&stats.TotalSent)
	if err != nil {
		return nil, fmt.Errorf("count sent: %w", err)
	}

	// Count total failed
	err = r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM notifications n
		JOIN users u ON n.user_id = u.id
		JOIN user_saved_places usp ON usp.user_id = u.id
		WHERE usp.place_id = $1 AND n.status = 'failed'
		AND n.created_at BETWEEN $2 AND $3
	`, placeID, from, to).Scan(&stats.TotalFailed)
	if err != nil {
		return nil, fmt.Errorf("count failed: %w", err)
	}

	// Count by channel
	rows, err := r.pool.Query(ctx, `
		SELECT n.channel, COUNT(*) FROM notifications n
		JOIN users u ON n.user_id = u.id
		JOIN user_saved_places usp ON usp.user_id = u.id
		WHERE usp.place_id = $1 AND n.status = 'sent'
		AND n.created_at BETWEEN $2 AND $3
		GROUP BY n.channel
	`, placeID, from, to)
	if err != nil {
		return nil, fmt.Errorf("count by channel: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var channel string
		var count int
		if err := rows.Scan(&channel, &count); err != nil {
			return nil, fmt.Errorf("scan channel count: %w", err)
		}
		stats.ByChannel[channel] = count
	}

	return stats, rows.Err()
}

// ReportStats queries community report metrics for an organization.
func (r *Repository) ReportStats(ctx context.Context, orgID uuid.UUID, from, to string) (*ReportStats, error) {
	stats := &ReportStats{
		OrgID:       orgID,
		PeriodStart: parseTime(from),
		PeriodEnd:   parseTime(to),
	}

	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM observations
		WHERE organization_id = $1 AND created_at BETWEEN $2 AND $3
	`, orgID, from, to).Scan(&stats.TotalSubmitted)
	if err != nil {
		return nil, fmt.Errorf("count submitted: %w", err)
	}

	err = r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM observations
		WHERE organization_id = $1 AND status = 'verified'
		AND created_at BETWEEN $2 AND $3
	`, orgID, from, to).Scan(&stats.TotalVerified)
	if err != nil {
		return nil, fmt.Errorf("count verified: %w", err)
	}

	err = r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM observations
		WHERE organization_id = $1 AND status = 'rejected'
		AND created_at BETWEEN $2 AND $3
	`, orgID, from, to).Scan(&stats.TotalRejected)
	if err != nil {
		return nil, fmt.Errorf("count rejected: %w", err)
	}

	err = r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM observations
		WHERE organization_id = $1 AND status = 'pending'
	`, orgID).Scan(&stats.PendingReview)
	if err != nil {
		return nil, fmt.Errorf("count pending: %w", err)
	}

	return stats, nil
}

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}
