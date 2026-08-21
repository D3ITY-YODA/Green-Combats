package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"green-compass-backend/pkg/database"
)

var ErrRefreshTokenNotFound = errors.New("refresh token not found")

type RefreshToken struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	TokenHash  string
	IssuedAt   time.Time
	ExpiresAt  time.Time
	RevokedAt  *time.Time
	ReplacedBy *uuid.UUID
}

type Repository struct {
	pool *database.Pool
}

func NewRepository(pool *database.Pool) *Repository {
	return &Repository{pool: pool}
}

const refreshColumns = `id, user_id, token_hash, issued_at, expires_at, revoked_at, replaced_by`

func (r *Repository) Create(ctx context.Context, rt *RefreshToken) error {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO refresh_tokens (user_id, token_hash, issued_at, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		rt.UserID, rt.TokenHash, rt.IssuedAt, rt.ExpiresAt)
	return row.Scan(&rt.ID)
}

func (r *Repository) ByHash(ctx context.Context, hash string) (*RefreshToken, error) {
	rt := &RefreshToken{}
	row := r.pool.QueryRow(ctx,
		fmt.Sprintf(`SELECT %s FROM refresh_tokens WHERE token_hash = $1`, refreshColumns), hash)
	err := row.Scan(&rt.ID, &rt.UserID, &rt.TokenHash, &rt.IssuedAt, &rt.ExpiresAt, &rt.RevokedAt, &rt.ReplacedBy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRefreshTokenNotFound
		}
		return nil, err
	}
	return rt, nil
}

func (r *Repository) Revoke(ctx context.Context, id uuid.UUID, replacedBy *uuid.UUID, now time.Time) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE refresh_tokens
		SET revoked_at = $2, replaced_by = $3
		WHERE id = $1 AND revoked_at IS NULL`,
		id, now, replacedBy)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrRefreshTokenNotFound
	}
	return nil
}

func (r *Repository) RevokeAllForUser(ctx context.Context, userID uuid.UUID, now time.Time) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE refresh_tokens
		SET revoked_at = $2
		WHERE user_id = $1 AND revoked_at IS NULL`,
		userID, now)
	return err
}
