package places

import (
	"context"
	"errors"
	"fmt"

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

const placeColumns = `id, name, place_type,
	ST_Y(location::geometry), ST_X(location::geometry),
	external_code, created_by, created_at, updated_at`

func (r *Repository) Create(ctx context.Context, p *Place) error {
	row := r.q.QueryRow(ctx, `
		INSERT INTO places (name, place_type, location, created_by)
		VALUES ($1, $2, ST_SetSRID(ST_MakePoint($3, $4), 4326)::geography, $5)
		RETURNING `+placeColumns,
		p.Name, p.PlaceType, p.Lon, p.Lat, p.CreatedBy)
	return scanPlace(row, p)
}

func (r *Repository) ByID(ctx context.Context, id uuid.UUID) (*Place, error) {
	p := &Place{}
	row := r.q.QueryRow(ctx,
		fmt.Sprintf(`SELECT %s FROM places WHERE id = $1`, placeColumns), id)
	if err := scanPlace(row, p); err != nil {
		return nil, mapReadError(err)
	}
	return p, nil
}

func (r *Repository) ByExternalCode(ctx context.Context, code string) (*Place, error) {
	p := &Place{}
	row := r.q.QueryRow(ctx,
		fmt.Sprintf(`SELECT %s FROM places WHERE external_code = $1`, placeColumns), code)
	if err := scanPlace(row, p); err != nil {
		return nil, mapReadError(err)
	}
	return p, nil
}

func (r *Repository) Update(ctx context.Context, p *Place) (*Place, error) {
	updated := &Place{}
	row := r.q.QueryRow(ctx, `
		UPDATE places SET name = $2, location = ST_SetSRID(ST_MakePoint($3, $4), 4326)::geography
		WHERE id = $1
		RETURNING `+placeColumns,
		p.ID, p.Name, p.Lon, p.Lat)
	if err := scanPlace(row, updated); err != nil {
		return nil, mapReadError(err)
	}
	return updated, nil
}

func (r *Repository) Nearby(ctx context.Context, lat, lon, radiusM float64, limit int) ([]WithDistance, error) {
	rows, err := r.q.Query(ctx, `
		SELECT `+placeColumns+`, ST_Distance(location, ref.point) AS distance_m
		FROM places,
			(SELECT ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography AS point) AS ref
		WHERE place_type IN ('community', 'ward', 'district')
			AND ST_DWithin(location, ref.point, $3)
		ORDER BY location <-> ref.point
		LIMIT $4`,
		lon, lat, radiusM, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []WithDistance
	for rows.Next() {
		var wd WithDistance
		if err := rows.Scan(&wd.ID, &wd.Name, &wd.PlaceType, &wd.Lat, &wd.Lon,
			&wd.ExternalCode, &wd.CreatedBy, &wd.CreatedAt, &wd.UpdatedAt, &wd.DistanceMeters); err != nil {
			return nil, err
		}
		result = append(result, wd)
	}
	return result, rows.Err()
}

func mapReadError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func scanPlace(row pgx.Row, p *Place) error {
	return row.Scan(&p.ID, &p.Name, &p.PlaceType, &p.Lat, &p.Lon,
		&p.ExternalCode, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
}
