package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"green-compass-backend/pkg/database"
)

type scanner interface {
	Scan(dest ...any) error
}

const userColumns = `id, phone_number, email, password_hash, display_name, language, is_platform_admin, created_at, updated_at`

type Repository struct {
	pool *database.Pool
}

func NewRepository(pool *database.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, u *User) error {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO users (phone_number, email, password_hash, display_name, language)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING `+userColumns,
		u.PhoneNumber, u.Email, u.PasswordHash, u.DisplayName, u.Language)
	return mapCreateError(scanUser(row, u))
}

func (r *Repository) ByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return r.by(ctx, "id", id)
}

func (r *Repository) ByEmail(ctx context.Context, email string) (*User, error) {
	return r.by(ctx, "email", email)
}

func (r *Repository) ByPhone(ctx context.Context, phone string) (*User, error) {
	return r.by(ctx, "phone_number", phone)
}

func (r *Repository) by(ctx context.Context, column string, value any) (*User, error) {
	u := &User{}
	row := r.pool.QueryRow(ctx,
		fmt.Sprintf(`SELECT %s FROM users WHERE %s = $1`, userColumns, column), value)
	if err := scanUser(row, u); err != nil {
		return nil, mapReadError(err)
	}
	return u, nil
}

func (r *Repository) UpdateProfile(ctx context.Context, id uuid.UUID, displayName, language string) (*User, error) {
	u := &User{}
	row := r.pool.QueryRow(ctx, `
		UPDATE users SET display_name = $2, language = $3
		WHERE id = $1
		RETURNING `+userColumns,
		id, displayName, language)
	if err := scanUser(row, u); err != nil {
		return nil, mapReadError(err)
	}
	return u, nil
}

func mapCreateError(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		switch pgErr.ConstraintName {
		case "users_email_key":
			return ErrEmailTaken
		case "users_phone_number_key":
			return ErrPhoneTaken
		}
	}
	return err
}

func mapReadError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func scanUser(row scanner, u *User) error {
	return row.Scan(&u.ID, &u.PhoneNumber, &u.Email, &u.PasswordHash,
		&u.DisplayName, &u.Language, &u.IsPlatformAdmin, &u.CreatedAt, &u.UpdatedAt)
}
