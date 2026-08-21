package organizations

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

const orgColumns = `id, name, org_type, status, created_at, updated_at`

func (r *Repository) Create(ctx context.Context, o *Organization) error {
	row := r.q.QueryRow(ctx, `
		INSERT INTO organizations (name, org_type)
		VALUES ($1, $2)
		RETURNING `+orgColumns,
		o.Name, o.OrgType)
	return mapWriteError(scanOrg(row, o))
}

func (r *Repository) ByID(ctx context.Context, id uuid.UUID) (*Organization, error) {
	o := &Organization{}
	row := r.q.QueryRow(ctx,
		fmt.Sprintf(`SELECT %s FROM organizations WHERE id = $1`, orgColumns), id)
	if err := scanOrg(row, o); err != nil {
		return nil, mapReadError(err)
	}
	return o, nil
}

func (r *Repository) SetStatus(ctx context.Context, id uuid.UUID, status string) (*Organization, error) {
	o := &Organization{}
	row := r.q.QueryRow(ctx, `
		UPDATE organizations SET status = $2
		WHERE id = $1
		RETURNING `+orgColumns,
		id, status)
	if err := scanOrg(row, o); err != nil {
		return nil, mapReadError(err)
	}
	return o, nil
}

func (r *Repository) AddMember(ctx context.Context, m *Member) error {
	tag, err := r.q.Exec(ctx, `
		INSERT INTO organization_members (organization_id, user_id, member_role)
		VALUES ($1, $2, $3)
		ON CONFLICT (organization_id, user_id) DO NOTHING`,
		m.OrgID, m.UserID, m.Role)
	if err != nil {
		return mapWriteError(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrAlreadyMember
	}
	return nil
}

func (r *Repository) RemoveMember(ctx context.Context, orgID, userID uuid.UUID) error {
	tag, err := r.q.Exec(ctx, `
		DELETE FROM organization_members
		WHERE organization_id = $1 AND user_id = $2`,
		orgID, userID)
	if err != nil {
		return mapWriteError(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotMember
	}
	return nil
}

func (r *Repository) MemberRole(ctx context.Context, orgID, userID uuid.UUID) (string, error) {
	var role string
	err := r.q.QueryRow(ctx, `
		SELECT member_role FROM organization_members
		WHERE organization_id = $1 AND user_id = $2`,
		orgID, userID).Scan(&role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotMember
		}
		return "", err
	}
	return role, nil
}

func mapReadError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func mapWriteError(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		switch pgErr.ConstraintName {
		case "organizations_name_unique":
			return ErrNameTaken
		case "organization_members_pkey":
			return ErrAlreadyMember
		}
	}
	return err
}

func scanOrg(row pgx.Row, o *Organization) error {
	return row.Scan(&o.ID, &o.Name, &o.OrgType, &o.Status, &o.CreatedAt, &o.UpdatedAt)
}
