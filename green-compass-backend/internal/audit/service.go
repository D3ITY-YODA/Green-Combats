package audit

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type repo interface {
	Create(ctx context.Context, e *Entry) error
	List(ctx context.Context, req ListRequest) ([]Entry, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Entry, error)
}

type Service struct {
	repo repo
}

func NewService(repo repo) *Service {
	return &Service{repo: repo}
}

// Record creates a new audit log entry.
func (s *Service) Record(ctx context.Context, req CreateRequest) (*Entry, error) {
	if req.Action == "" {
		return nil, fmt.Errorf("action is required")
	}
	if req.ResourceType == "" {
		return nil, fmt.Errorf("resource_type is required")
	}

	e := &Entry{
		ActorID:      req.ActorID,
		ActorOrgID:   req.ActorOrgID,
		Action:       req.Action,
		ResourceType: req.ResourceType,
		ResourceID:   req.ResourceID,
		Details:      req.Details,
		IPAddress:    req.IPAddress,
	}

	if err := s.repo.Create(ctx, e); err != nil {
		return nil, fmt.Errorf("create audit entry: %w", err)
	}
	return e, nil
}

// List returns paginated audit entries matching the given filters.
func (s *Service) List(ctx context.Context, req ListRequest) (*ListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 20
	}

	entries, total, err := s.repo.List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("list audit entries: %w", err)
	}

	return &ListResponse{
		Entries: entries,
		Page:    req.Page,
		Limit:   req.Limit,
		Total:   total,
	}, nil
}

// GetByID returns a single audit entry.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Entry, error) {
	return s.repo.GetByID(ctx, id)
}
