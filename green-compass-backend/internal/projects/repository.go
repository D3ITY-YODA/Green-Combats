package projects

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("project not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, p *Project) error {
	query := `
		INSERT INTO projects (org_id, name, description, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	err := r.pool.QueryRow(ctx, query,
		p.OrgID, p.Name, p.Description, p.StartDate, p.EndDate,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert project: %w", err)
	}
	p.Status = "active"
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Project, error) {
	query := `
		SELECT id, org_id, name, description, status, start_date, end_date, created_at, updated_at
		FROM projects WHERE id = $1
	`
	var p Project
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.OrgID, &p.Name, &p.Description, &p.Status,
		&p.StartDate, &p.EndDate, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}
	return &p, nil
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*Project, error) {
	p, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Description != nil {
		p.Description = *req.Description
	}
	if req.Status != nil {
		p.Status = *req.Status
	}
	if req.StartDate != nil {
		p.StartDate = req.StartDate
	}
	if req.EndDate != nil {
		p.EndDate = req.EndDate
	}

	query := `
		UPDATE projects SET name = $2, description = $3, status = $4, start_date = $5, end_date = $6
		WHERE id = $1
		RETURNING updated_at
	`
	err = r.pool.QueryRow(ctx, query, id, p.Name, p.Description, p.Status, p.StartDate, p.EndDate).Scan(&p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("update project: %w", err)
	}
	return p, nil
}

func (r *Repository) List(ctx context.Context, req ListRequest) ([]Project, int, error) {
	countQuery := `SELECT COUNT(*) FROM projects WHERE org_id = $1`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, req.OrgID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count projects: %w", err)
	}

	query := `
		SELECT id, org_id, name, description, status, start_date, end_date, created_at, updated_at
		FROM projects WHERE org_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.pool.Query(ctx, query, req.OrgID, req.Limit, (req.Page-1)*req.Limit)
	if err != nil {
		return nil, 0, fmt.Errorf("query projects: %w", err)
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(
			&p.ID, &p.OrgID, &p.Name, &p.Description, &p.Status,
			&p.StartDate, &p.EndDate, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan project: %w", err)
		}
		projects = append(projects, p)
	}
	return projects, total, rows.Err()
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM projects WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete project: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
