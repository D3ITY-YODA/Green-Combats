package observations

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"green-compass-backend/pkg/database"
)

type querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type Repository struct {
	q querier
}

func NewRepository(pool *database.Pool) *Repository {
	return &Repository{q: pool}
}

func (r *Repository) WithTx(tx pgx.Tx) *Repository {
	return &Repository{q: tx}
}

const obsColumns = `id, reporter_id, place_id,
	ST_Y(location::geometry), ST_X(location::geometry),
	category, description, photo_object_key, status,
	verified_by, verified_at, created_at, updated_at`

func (r *Repository) Create(ctx context.Context, o *Observation) error {
	var lat, lon *float64
	if o.Lat != nil && o.Lon != nil {
		lat = o.Lat
		lon = o.Lon
	}
	row := r.q.QueryRow(ctx, `
		INSERT INTO observations (reporter_id, place_id, location, category, description, photo_object_key)
		VALUES ($1, $2, ST_SetSRID(ST_MakePoint($3, $4), 4326)::geography, $5, $6, $7)
		RETURNING `+obsColumns,
		o.ReporterID, o.PlaceID, lon, lat, o.Category, o.Description, o.PhotoObjectKey)
	return scanObservation(row, o)
}

func (r *Repository) ByID(ctx context.Context, id uuid.UUID) (*Observation, error) {
	o := &Observation{}
	row := r.q.QueryRow(ctx,
		fmt.Sprintf(`SELECT %s FROM observations WHERE id = $1`, obsColumns), id)
	if err := scanObservation(row, o); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return o, nil
}

func (r *Repository) ListByReporter(ctx context.Context, reporterID uuid.UUID, limit int) ([]Observation, error) {
	rows, err := r.q.Query(ctx,
		fmt.Sprintf(`SELECT %s FROM observations WHERE reporter_id = $1 ORDER BY created_at DESC LIMIT $2`, obsColumns),
		reporterID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanObservations(rows)
}

func (r *Repository) ListPending(ctx context.Context, limit int) ([]Observation, error) {
	rows, err := r.q.Query(ctx,
		fmt.Sprintf(`SELECT %s FROM observations WHERE status = $1 ORDER BY created_at DESC LIMIT $2`, obsColumns),
		StatusPending, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanObservations(rows)
}

func (r *Repository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, verifiedBy *uuid.UUID, verifiedAt *time.Time) (*Observation, error) {
	o := &Observation{}
	row := r.q.QueryRow(ctx, `
		UPDATE observations SET status = $2, verified_by = $3, verified_at = $4
		WHERE id = $1
		RETURNING `+obsColumns,
		id, status, verifiedBy, verifiedAt)
	if err := scanObservation(row, o); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return o, nil
}

func scanObservation(row pgx.Row, o *Observation) error {
	return row.Scan(&o.ID, &o.ReporterID, &o.PlaceID, &o.Lat, &o.Lon,
		&o.Category, &o.Description, &o.PhotoObjectKey, &o.Status,
		&o.VerifiedBy, &o.VerifiedAt, &o.CreatedAt, &o.UpdatedAt)
}

func scanObservations(rows pgx.Rows) ([]Observation, error) {
	var result []Observation
	for rows.Next() {
		var o Observation
		if err := rows.Scan(&o.ID, &o.ReporterID, &o.PlaceID, &o.Lat, &o.Lon,
			&o.Category, &o.Description, &o.PhotoObjectKey, &o.Status,
			&o.VerifiedBy, &o.VerifiedAt, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, o)
	}
	return result, rows.Err()
}

func mapWriteError(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		switch pgErr.ConstraintName {
		case "observations_place_id_fkey":
			return fmt.Errorf("place not found")
		}
	}
	return err
}
