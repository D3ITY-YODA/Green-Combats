package outbox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("outbox event not found")

// Event is a transactional outbox row published asynchronously.
type Event struct {
	ID            uuid.UUID
	EventType     string
	AggregateType string
	AggregateID   uuid.UUID
	Payload       json.RawMessage
	Status        string // pending, published, failed
	Attempts      int
	AvailableAt   time.Time
	PublishedAt   *time.Time
	CreatedAt     time.Time
}

const (
	StatusPending   = "pending"
	StatusPublished = "published"
	StatusFailed    = "failed"

	maxAttempts = 5
)

// Repository persists and polls outbox events.
type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Publish stores an event in the same transaction as the caller's work
// when tx is non-nil; otherwise it commits on its own.
func (r *Repository) Publish(ctx context.Context, tx pgx.Tx, eventType, aggregateType string, aggregateID uuid.UUID, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal outbox payload: %w", err)
	}

	query := `
		INSERT INTO outbox_events (event_type, aggregate_type, aggregate_id, payload, status, attempts, available_at)
		VALUES ($1, $2, $3, $4, $5, 0, now())
	`

	if tx != nil {
		_, err = tx.Exec(ctx, query, eventType, aggregateType, aggregateID, data, StatusPending)
	} else {
		_, err = r.pool.Exec(ctx, query, eventType, aggregateType, aggregateID, data, StatusPending)
	}
	if err != nil {
		return fmt.Errorf("insert outbox event: %w", err)
	}
	return nil
}

// PollBatch claims up to limit pending events whose availability time has
// passed, marking them as being processed via FOR UPDATE SKIP LOCKED.
func (r *Repository) PollBatch(ctx context.Context, limit int) ([]Event, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin poll tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	rows, err := tx.Query(ctx, `
		SELECT id, event_type, aggregate_type, aggregate_id, payload, status, attempts, available_at, created_at
		FROM outbox_events
		WHERE status = $1 AND attempts < $2 AND available_at <= now()
		ORDER BY created_at
		LIMIT $3
		FOR UPDATE SKIP LOCKED
	`, StatusPending, maxAttempts, limit)
	if err != nil {
		return nil, fmt.Errorf("poll outbox: %w", err)
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.EventType, &e.AggregateType, &e.AggregateID,
			&e.Payload, &e.Status, &e.Attempts, &e.AvailableAt, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan outbox event: %w", err)
		}
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, tx.Commit(ctx)
}

// MarkPublished records successful delivery.
func (r *Repository) MarkPublished(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE outbox_events SET status = $1, published_at = now() WHERE id = $2
	`, StatusPublished, id)
	return err
}

// MarkFailed increments the attempt counter and backs off availability,
// or moves the event to failed once retries are exhausted.
func (r *Repository) MarkFailed(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE outbox_events
		SET attempts = attempts + 1,
		    status = CASE WHEN attempts + 1 >= $1 THEN $2 ELSE status END,
		    available_at = now() + make_interval(secs => power(2, attempts + 1))
		WHERE id = $3
	`, maxAttempts, StatusFailed, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// Dispatcher drains pending events through a handler with retry/backoff.
type Dispatcher struct {
	repo   *Repository
	handle func(ctx context.Context, e Event) error
}

func NewDispatcher(repo *Repository, handle func(ctx context.Context, e Event) error) *Dispatcher {
	return &Dispatcher{repo: repo, handle: handle}
}

// ProcessOnce handles one batch; returns number processed. Intended to be
// called on a ticker from cmd/worker or cmd/notifier.
func (d *Dispatcher) ProcessOnce(ctx context.Context, batchSize int) (int, error) {
	events, err := d.repo.PollBatch(ctx, batchSize)
	if err != nil {
		return 0, err
	}
	for _, e := range events {
		if err := d.handle(ctx, e); err != nil {
			_ = d.repo.MarkFailed(ctx, e.ID)
			continue
		}
		if err := d.repo.MarkPublished(ctx, e.ID); err != nil {
			return len(events), err
		}
	}
	return len(events), nil
}
