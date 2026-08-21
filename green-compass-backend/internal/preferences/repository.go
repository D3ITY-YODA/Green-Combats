package preferences

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"green-compass-backend/pkg/database"
)

type Repository struct {
	pool *database.Pool
}

func NewRepository(pool *database.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) List(ctx context.Context, userID uuid.UUID) ([]SavedPlace, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT p.id, p.name, p.place_type, ST_Y(p.location::geometry), ST_X(p.location::geometry),
			p.external_code, p.created_by, p.created_at, p.updated_at,
			usp.label, usp.is_primary, usp.saved_at
		FROM user_saved_places usp
		JOIN places p ON p.id = usp.place_id
		WHERE usp.user_id = $1
		ORDER BY usp.is_primary DESC, usp.saved_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []SavedPlace
	for rows.Next() {
		var saved SavedPlace
		if err := rows.Scan(&saved.ID, &saved.Name, &saved.PlaceType, &saved.Lat, &saved.Lon,
			&saved.ExternalCode, &saved.CreatedBy, &saved.CreatedAt, &saved.UpdatedAt,
			&saved.Label, &saved.IsPrimary, &saved.SavedAt); err != nil {
			return nil, err
		}
		result = append(result, saved)
	}
	return result, rows.Err()
}

func (r *Repository) Save(ctx context.Context, userID, placeID uuid.UUID, label *string) error {
	return database.WithTx(ctx, r.pool, func(ctx context.Context, tx pgx.Tx) error {
		if err := lockUser(ctx, tx, userID); err != nil {
			return err
		}
		var exists bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM user_saved_places WHERE user_id = $1 AND place_id = $2)`, userID, placeID).Scan(&exists); err != nil {
			return err
		}
		if !exists {
			var count int
			if err := tx.QueryRow(ctx, `SELECT count(*) FROM user_saved_places WHERE user_id = $1`, userID).Scan(&count); err != nil {
				return err
			}
			if count >= MaxSavedPlaces {
				return ErrSavedPlaceLimit
			}
		}
		_, err := tx.Exec(ctx, `
			INSERT INTO user_saved_places (user_id, place_id, label)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id, place_id) DO UPDATE SET label = EXCLUDED.label`, userID, placeID, label)
		return mapSaveError(err)
	})
}

func (r *Repository) Unsave(ctx context.Context, userID, placeID uuid.UUID) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM user_saved_places WHERE user_id = $1 AND place_id = $2`, userID, placeID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotSaved
	}
	return nil
}

func (r *Repository) SetPrimary(ctx context.Context, userID, placeID uuid.UUID) error {
	return database.WithTx(ctx, r.pool, func(ctx context.Context, tx pgx.Tx) error {
		if err := lockUser(ctx, tx, userID); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `UPDATE user_saved_places SET is_primary = FALSE WHERE user_id = $1 AND is_primary`, userID); err != nil {
			return err
		}
		result, err := tx.Exec(ctx, `UPDATE user_saved_places SET is_primary = TRUE WHERE user_id = $1 AND place_id = $2`, userID, placeID)
		if err != nil {
			return err
		}
		if result.RowsAffected() == 0 {
			return ErrNotSaved
		}
		return nil
	})
}

func lockUser(ctx context.Context, tx pgx.Tx, userID uuid.UUID) error {
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `SELECT id FROM users WHERE id = $1 FOR UPDATE`, userID).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("user not found")
		}
		return err
	}
	return nil
}

func mapSaveError(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" && pgErr.ConstraintName == "user_saved_places_place_id_fkey" {
		return ErrPlaceNotFound
	}
	return err
}
