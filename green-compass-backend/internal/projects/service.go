package projects

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type repo interface {
	Create(ctx context.Context, p *Project) error
	GetByID(ctx context.Context, id uuid.UUID) (*Project, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*Project, error)
	List(ctx context.Context, req ListRequest) ([]Project, int, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	repo repo
}

func NewService(repo repo) *Service {
	return &Service{repo: repo}
}

// Create creates a new project for an organization.
func (s *Service) Create(ctx context.Context, req CreateRequest) (*Project, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	p := &Project{
		OrgID:       req.OrgID,
		Name:        req.Name,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}
	return p, nil
}

// GetByID returns a single project.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Project, error) {
	return s.repo.GetByID(ctx, id)
}

// Update modifies an existing project.
func (s *Service) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*Project, error) {
	return s.repo.Update(ctx, id, req)
}

// List returns paginated projects for an organization.
func (s *Service) List(ctx context.Context, req ListRequest) (*ListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 20
	}

	projects, total, err := s.repo.List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}

	return &ListResponse{
		Projects: projects,
		Page:     req.Page,
		Limit:    req.Limit,
		Total:    total,
	}, nil
}

// Delete removes a project.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
