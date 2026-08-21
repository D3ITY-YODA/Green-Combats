package updates

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"green-compass-backend/pkg/httpx"
)

type updatesRepo interface {
	GetTodayForPlace(ctx context.Context, placeID uuid.UUID) (*Update, error)
	ListForUser(ctx context.Context, req ListRequest) ([]Update, int, error)
	ListIndicators(ctx context.Context, contentID uuid.UUID) ([]Indicator, error)
	ExploreIndicators(ctx context.Context, req ExploreRequest) ([]ExploreIndicator, int, error)
}

type Service struct {
	repo   updatesRepo
	logger *slog.Logger
}

func NewService(repo updatesRepo, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// GetToday retrieves today's update for user's current place
func (s *Service) GetToday(ctx context.Context, placeID uuid.UUID) (*Update, error) {
	update, err := s.repo.GetTodayForPlace(ctx, placeID)
	if err != nil {
		return nil, err
	}

	// Fetch indicators for this update
	indicators, err := s.repo.ListIndicators(ctx, update.ID)
	if err != nil {
		s.logger.Warn("failed to fetch indicators for update", "update_id", update.ID, "err", err)
		// Don't fail if indicators unavailable, just return update without them
	} else {
		update.Indicators = indicators
	}

	return update, nil
}

// List retrieves paginated updates for user's places
func (s *Service) List(ctx context.Context, req ListRequest) (*ListResponse, error) {
	// Validate pagination
	paging := httpx.ValidatePagination(req.Page, req.Limit)
	req.Page = paging.Page
	req.Limit = paging.Limit

	updates, total, err := s.repo.ListForUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return &ListResponse{
		Updates: updates,
		Total:   total,
		Page:    req.Page,
		Limit:   req.Limit,
		HasNext: req.Page < (total + req.Limit - 1) / req.Limit,
	}, nil
}

// Explore retrieves indicators for discovery browsing
func (s *Service) Explore(ctx context.Context, req ExploreRequest) (*ExploreResponse, error) {
	// Validate pagination
	paging := httpx.ValidatePagination(req.Page, req.Limit)
	req.Page = paging.Page
	req.Limit = paging.Limit

	indicators, total, err := s.repo.ExploreIndicators(ctx, req)
	if err != nil {
		return nil, err
	}

	return &ExploreResponse{
		Indicators: indicators,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		HasNext:    req.Page < (total + req.Limit - 1) / req.Limit,
	}, nil
}
