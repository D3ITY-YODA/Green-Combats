package ingestion

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"green-compass-backend/pkg/database"
)

type Repository struct{ pool *database.Pool }

func NewRepository(pool *database.Pool) *Repository { return &Repository{pool: pool} }

func (r *Repository) StartRun(ctx context.Context, sourceID, placeID uuid.UUID, from, to *time.Time) (*Run, error) {
	run := &Run{}
	err := r.pool.QueryRow(ctx, `
		INSERT INTO ingestion_runs (source_id, place_id, status, requested_from, requested_to)
		VALUES ($1, $2, 'running', $3, $4)
		RETURNING id, source_id, place_id, status, requested_from, requested_to, started_at, finished_at, landed_record_count, error_message`,
		sourceID, placeID, from, to).Scan(&run.ID, &run.SourceID, &run.PlaceID, &run.Status, &run.RequestedFrom, &run.RequestedTo, &run.StartedAt, &run.FinishedAt, &run.LandedRecordCount, &run.ErrorMessage)
	if err != nil {
		return nil, err
	}
	return run, nil
}

// StoreRaw inserts an immutable provider payload. inserted is false when the
// same source/place/observation-time/checksum has already been stored.
func (r *Repository) StoreRaw(ctx context.Context, record *RawRecord) (inserted bool, err error) {
	if !json.Valid(record.Payload) {
		return false, fmt.Errorf("%w: raw payload must be valid JSON", ErrInvalidData)
	}
	if record.PayloadChecksum == "" {
		record.PayloadChecksum = Checksum(record.Payload)
	}
	err = r.pool.QueryRow(ctx, `
		INSERT INTO raw_records (
			ingestion_run_id, source_id, place_id, source_observed_at,
			source_url, content_type, payload, payload_checksum
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (source_id, place_id, source_observed_at, payload_checksum) DO NOTHING
		RETURNING id, fetched_at, created_at`,
		record.IngestionRunID, record.SourceID, record.PlaceID, record.SourceObservedAt,
		record.SourceURL, record.ContentType, record.Payload, record.PayloadChecksum).Scan(&record.ID, &record.FetchedAt, &record.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Repository) FinishRun(ctx context.Context, id uuid.UUID, status string, landedCount int, message *string) error {
	if status != StatusSucceeded && status != StatusFailed {
		return fmt.Errorf("%w: invalid final run status", ErrInvalidData)
	}
	result, err := r.pool.Exec(ctx, `
		UPDATE ingestion_runs
		SET status = $2, finished_at = now(), landed_record_count = $3, error_message = $4
		WHERE id = $1 AND status = 'running'`, id, status, landedCount, message)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return fmt.Errorf("%w: ingestion run is not running", ErrInvalidData)
	}
	return nil
}
