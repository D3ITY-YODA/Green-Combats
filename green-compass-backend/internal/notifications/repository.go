package notifications

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("notification not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// --- Preferences ---

func (r *Repository) UpsertPreference(ctx context.Context, p *Preference) error {
	query := `
		INSERT INTO notification_preferences (user_id, channel, event_type, enabled)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, channel, event_type) DO UPDATE SET enabled = $4
	`
	_, err := r.pool.Exec(ctx, query, p.UserID, p.Channel, p.EventType, p.Enabled)
	if err != nil {
		return fmt.Errorf("upsert preference: %w", err)
	}
	return nil
}

func (r *Repository) GetPreferences(ctx context.Context, userID uuid.UUID) ([]Preference, error) {
	query := `
		SELECT user_id, channel, event_type, enabled, created_at, updated_at
		FROM notification_preferences
		WHERE user_id = $1
		ORDER BY channel, event_type
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query preferences: %w", err)
	}
	defer rows.Close()

	var prefs []Preference
	for rows.Next() {
		var p Preference
		if err := rows.Scan(&p.UserID, &p.Channel, &p.EventType, &p.Enabled, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan preference: %w", err)
		}
		prefs = append(prefs, p)
	}
	return prefs, rows.Err()
}

// IsEnabled checks if a user has notifications enabled for a given channel and event type.
// Defaults to true if no preference record exists (opt-out model).
func (r *Repository) IsEnabled(ctx context.Context, userID uuid.UUID, channel Channel, eventType string) (bool, error) {
	query := `
		SELECT enabled
		FROM notification_preferences
		WHERE user_id = $1 AND channel = $2 AND event_type = $3
	`
	var enabled bool
	err := r.pool.QueryRow(ctx, query, userID, channel, eventType).Scan(&enabled)
	if errors.Is(err, pgx.ErrNoRows) {
		return true, nil // default: opt-out model
	}
	if err != nil {
		return false, fmt.Errorf("check preference: %w", err)
	}
	return enabled, nil
}

// --- Notifications ---

func (r *Repository) Create(ctx context.Context, n *Notification) error {
	query := `
		INSERT INTO notifications (user_id, channel, event_type, title, body, metadata, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`
	err := r.pool.QueryRow(ctx, query,
		n.UserID, n.Channel, n.EventType, n.Title, n.Body,
		n.Metadata, n.Status,
	).Scan(&n.ID, &n.CreatedAt)
	if err != nil {
		return fmt.Errorf("create notification: %w", err)
	}
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Notification, error) {
	query := `
		SELECT id, user_id, channel, event_type, title, body, metadata, status, created_at, sent_at, error_message
		FROM notifications
		WHERE id = $1
	`
	var n Notification
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&n.ID, &n.UserID, &n.Channel, &n.EventType, &n.Title, &n.Body,
		&n.Metadata, &n.Status, &n.CreatedAt, &n.SentAt, &n.ErrorMessage,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get notification: %w", err)
	}
	return &n, nil
}

func (r *Repository) List(ctx context.Context, req ListRequest) ([]Notification, int, error) {
	countQuery := `SELECT COUNT(*) FROM notifications WHERE user_id = $1`
	args := []interface{}{req.UserID}

	if req.Status != nil {
		countQuery += ` AND status = $2`
		args = append(args, *req.Status)
	}

	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count notifications: %w", err)
	}

	query := `
		SELECT id, user_id, channel, event_type, title, body, metadata, status, created_at, sent_at, error_message
		FROM notifications
		WHERE user_id = $1
	`
	queryArgs := []interface{}{req.UserID}
	paramIdx := 2
	if req.Status != nil {
		query += fmt.Sprintf(` AND status = $%d`, paramIdx)
		queryArgs = append(queryArgs, *req.Status)
		paramIdx++
	}

	query += ` ORDER BY created_at DESC`
	query += fmt.Sprintf(` LIMIT $%d`, paramIdx)
	queryArgs = append(queryArgs, req.Limit)
	paramIdx++
	query += fmt.Sprintf(` OFFSET $%d`, paramIdx)
	queryArgs = append(queryArgs, (req.Page-1)*req.Limit)

	rows, err := r.pool.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("list notifications: %w", err)
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var n Notification
		if err := rows.Scan(
			&n.ID, &n.UserID, &n.Channel, &n.EventType, &n.Title, &n.Body,
			&n.Metadata, &n.Status, &n.CreatedAt, &n.SentAt, &n.ErrorMessage,
		); err != nil {
			return nil, 0, fmt.Errorf("scan notification: %w", err)
		}
		notifications = append(notifications, n)
	}
	return notifications, total, rows.Err()
}

func (r *Repository) ListPending(ctx context.Context, limit int) ([]Notification, error) {
	query := `
		SELECT id, user_id, channel, event_type, title, body, metadata, status, created_at, sent_at, error_message
		FROM notifications
		WHERE status = 'pending'
		ORDER BY created_at ASC
		LIMIT $1
	`
	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("query pending: %w", err)
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var n Notification
		if err := rows.Scan(
			&n.ID, &n.UserID, &n.Channel, &n.EventType, &n.Title, &n.Body,
			&n.Metadata, &n.Status, &n.CreatedAt, &n.SentAt, &n.ErrorMessage,
		); err != nil {
			return nil, fmt.Errorf("scan notification: %w", err)
		}
		notifications = append(notifications, n)
	}
	return notifications, rows.Err()
}

func (r *Repository) MarkSent(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE notifications SET status = 'sent', sent_at = now() WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("mark sent: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) MarkFailed(ctx context.Context, id uuid.UUID, errMsg string) error {
	query := `UPDATE notifications SET status = 'failed', error_message = $2 WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id, errMsg)
	if err != nil {
		return fmt.Errorf("mark failed: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// --- Delivery Log ---

func (r *Repository) LogDelivery(ctx context.Context, log *DeliveryLog) error {
	query := `
		INSERT INTO notification_delivery_log (notification_id, attempt, status, provider_response)
		VALUES ($1, $2, $3, $4)
		RETURNING id, attempted_at
	`
	err := r.pool.QueryRow(ctx, query,
		log.NotificationID, log.Attempt, log.Status, log.ProviderResponse,
	).Scan(&log.ID, &log.AttemptedAt)
	if err != nil {
		return fmt.Errorf("log delivery: %w", err)
	}
	return nil
}
