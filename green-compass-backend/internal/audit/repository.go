package audit

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, e *Entry) error {
	query := `
		INSERT INTO audit_log (actor_id, actor_org_id, action, resource_type, resource_id, details, ip_address)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`
	err := r.pool.QueryRow(ctx, query,
		e.ActorID, e.ActorOrgID, e.Action, e.ResourceType,
		e.ResourceID, e.Details, e.IPAddress,
	).Scan(&e.ID, &e.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert audit entry: %w", err)
	}
	return nil
}

func (r *Repository) List(ctx context.Context, req ListRequest) ([]Entry, int, error) {
	countQuery := `SELECT COUNT(*) FROM audit_log WHERE 1=1`
	args := []interface{}{}
	paramIdx := 1

	if req.ActorID != nil {
		countQuery += fmt.Sprintf(` AND actor_id = $%d`, paramIdx)
		args = append(args, *req.ActorID)
		paramIdx++
	}
	if req.ActorOrgID != nil {
		countQuery += fmt.Sprintf(` AND actor_org_id = $%d`, paramIdx)
		args = append(args, *req.ActorOrgID)
		paramIdx++
	}
	if req.Action != nil {
		countQuery += fmt.Sprintf(` AND action = $%d`, paramIdx)
		args = append(args, *req.Action)
		paramIdx++
	}
	if req.ResourceType != nil {
		countQuery += fmt.Sprintf(` AND resource_type = $%d`, paramIdx)
		args = append(args, *req.ResourceType)
		paramIdx++
	}
	if req.ResourceID != nil {
		countQuery += fmt.Sprintf(` AND resource_id = $%d`, paramIdx)
		args = append(args, *req.ResourceID)
		paramIdx++
	}

	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count audit entries: %w", err)
	}

	query := `
		SELECT id, actor_id, actor_org_id, action, resource_type, resource_id, details, ip_address, created_at
		FROM audit_log WHERE 1=1
	`
	queryArgs := []interface{}{}
	queryParamIdx := 1

	if req.ActorID != nil {
		query += fmt.Sprintf(` AND actor_id = $%d`, queryParamIdx)
		queryArgs = append(queryArgs, *req.ActorID)
		queryParamIdx++
	}
	if req.ActorOrgID != nil {
		query += fmt.Sprintf(` AND actor_org_id = $%d`, queryParamIdx)
		queryArgs = append(queryArgs, *req.ActorOrgID)
		queryParamIdx++
	}
	if req.Action != nil {
		query += fmt.Sprintf(` AND action = $%d`, queryParamIdx)
		queryArgs = append(queryArgs, *req.Action)
		queryParamIdx++
	}
	if req.ResourceType != nil {
		query += fmt.Sprintf(` AND resource_type = $%d`, queryParamIdx)
		queryArgs = append(queryArgs, *req.ResourceType)
		queryParamIdx++
	}
	if req.ResourceID != nil {
		query += fmt.Sprintf(` AND resource_id = $%d`, queryParamIdx)
		queryArgs = append(queryArgs, *req.ResourceID)
		queryParamIdx++
	}

	query += ` ORDER BY created_at DESC`
	query += fmt.Sprintf(` LIMIT $%d`, queryParamIdx)
	queryArgs = append(queryArgs, req.Limit)
	queryParamIdx++
	query += fmt.Sprintf(` OFFSET $%d`, queryParamIdx)
	queryArgs = append(queryArgs, (req.Page-1)*req.Limit)

	rows, err := r.pool.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("query audit entries: %w", err)
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var e Entry
		if err := rows.Scan(
			&e.ID, &e.ActorID, &e.ActorOrgID, &e.Action, &e.ResourceType,
			&e.ResourceID, &e.Details, &e.IPAddress, &e.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan audit entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, total, rows.Err()
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Entry, error) {
	query := `
		SELECT id, actor_id, actor_org_id, action, resource_type, resource_id, details, ip_address, created_at
		FROM audit_log WHERE id = $1
	`
	var e Entry
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&e.ID, &e.ActorID, &e.ActorOrgID, &e.Action, &e.ResourceType,
		&e.ResourceID, &e.Details, &e.IPAddress, &e.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("audit entry not found")
		}
		return nil, fmt.Errorf("get audit entry: %w", err)
	}
	return &e, nil
}
