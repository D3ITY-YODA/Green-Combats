package reports

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"green-compass-backend/internal/observations"
)

type reportsRepo interface {
	ListPending(ctx context.Context, limit int) ([]Report, error)
	ByID(ctx context.Context, id uuid.UUID) (*Report, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, verifiedBy uuid.UUID) (*Report, error)
}

type Service struct {
	repo   reportsRepo
	logger *slog.Logger
}

func NewService(repo reportsRepo, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

func (s *Service) ListPending(ctx context.Context, caller Caller, limit int) ([]Report, error) {
	if !caller.IsPlatformAdmin {
		return nil, ErrNotAllowed
	}
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	return s.repo.ListPending(ctx, limit)
}

func (s *Service) Get(ctx context.Context, id uuid.UUID, caller Caller) (*Report, error) {
	if !caller.IsPlatformAdmin {
		return nil, ErrNotAllowed
	}
	return s.repo.ByID(ctx, id)
}

func (s *Service) Verify(ctx context.Context, id uuid.UUID, caller Caller, status string) (*Report, error) {
	if !caller.IsPlatformAdmin {
		return nil, ErrNotAllowed
	}
	if status != observations.StatusVerified && status != observations.StatusRejected && status != observations.StatusFlagged {
		return nil, fmt.Errorf("%w: invalid status %q", ErrInvalidData, status)
	}
	return s.repo.UpdateStatus(ctx, id, status, caller.UserID)
}
