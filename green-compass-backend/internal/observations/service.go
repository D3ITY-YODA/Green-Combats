package observations

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"

	"green-compass-backend/pkg/storage"
)

type observationRepo interface {
	Create(ctx context.Context, o *Observation) error
	ByID(ctx context.Context, id uuid.UUID) (*Observation, error)
	ListByReporter(ctx context.Context, reporterID uuid.UUID, limit int) ([]Observation, error)
	ListPending(ctx context.Context, limit int) ([]Observation, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, verifiedBy *uuid.UUID, verifiedAt *time.Time) (*Observation, error)
}

type Service struct {
	repo    observationRepo
	storage storage.Store
	logger  *slog.Logger
}

func NewService(repo observationRepo, store storage.Store, logger *slog.Logger) *Service {
	return &Service{
		repo:    repo,
		storage: store,
		logger:  logger,
	}
}

type CreateInput struct {
	PlaceID        *uuid.UUID
	Lat            *float64
	Lon            *float64
	Category       string
	Description    string
	PhotoObjectKey *string
	Caller         Caller
}

func (s *Service) Create(ctx context.Context, in CreateInput) (*Observation, error) {
	if !IsValidCategory(in.Category) {
		return nil, fmt.Errorf("%w: unknown category %q", ErrInvalidData, in.Category)
	}
	description := strings.TrimSpace(in.Description)
	if description == "" {
		return nil, fmt.Errorf("%w: description is required", ErrInvalidData)
	}
	if len(description) > 2000 {
		return nil, fmt.Errorf("%w: description exceeds 2000 characters", ErrInvalidData)
	}
	if in.PlaceID == nil && (in.Lat == nil || in.Lon == nil) {
		return nil, fmt.Errorf("%w: either place_id or lat/lon is required", ErrInvalidData)
	}

	o := &Observation{
		ReporterID:     in.Caller.UserID,
		PlaceID:        in.PlaceID,
		Lat:            in.Lat,
		Lon:            in.Lon,
		Category:       in.Category,
		Description:    description,
		PhotoObjectKey: in.PhotoObjectKey,
		Status:         StatusPending,
	}
	if err := s.repo.Create(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *Service) ByID(ctx context.Context, id uuid.UUID, caller Caller) (*Observation, error) {
	o, err := s.repo.ByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !canView(o, caller) {
		return nil, ErrNotAllowed
	}
	return o, nil
}

func (s *Service) ListMine(ctx context.Context, userID uuid.UUID, limit int) ([]Observation, error) {
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	return s.repo.ListByReporter(ctx, userID, limit)
}

func (s *Service) ListPending(ctx context.Context, caller Caller, limit int) ([]Observation, error) {
	if !caller.IsPlatformAdmin {
		return nil, ErrNotAllowed
	}
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	return s.repo.ListPending(ctx, limit)
}

func (s *Service) Verify(ctx context.Context, id uuid.UUID, caller Caller, status string) (*Observation, error) {
	if !caller.IsPlatformAdmin {
		return nil, ErrNotAllowed
	}
	if !IsValidStatus(status) || status == StatusPending {
		return nil, fmt.Errorf("%w: invalid status %q", ErrInvalidData, status)
	}
	now := time.Now().UTC()
	o, err := s.repo.UpdateStatus(ctx, id, status, &caller.UserID, &now)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func canView(o *Observation, caller Caller) bool {
	if caller.IsPlatformAdmin {
		return true
	}
	return o.ReporterID == caller.UserID
}
