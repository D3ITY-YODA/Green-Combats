package context

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"green-compass-backend/internal/places"
	"green-compass-backend/internal/preferences"
	"green-compass-backend/internal/updates"
)

type todayGetter interface {
	GetToday(ctx context.Context, placeID uuid.UUID) (*updates.Update, error)
}

type placeLister interface {
	List(ctx context.Context, userID uuid.UUID) ([]preferences.SavedPlace, error)
}

type placeGetter interface {
	ByID(ctx context.Context, id uuid.UUID) (*places.Place, error)
}

type Service struct {
	prefs  placeLister
	places placeGetter
	today  todayGetter
	logger *slog.Logger
}

func NewService(prefs placeLister, places placeGetter, today todayGetter, logger *slog.Logger) *Service {
	return &Service{
		prefs:  prefs,
		places: places,
		today:  today,
		logger: logger,
	}
}

func (s *Service) ResolveContext(ctx context.Context, req ResolveContextRequest) (*ContextResponse, error) {
	placeInfo, err := s.resolvePlace(ctx, req.UserID, req.PlaceID)
	if err != nil {
		return nil, err
	}

	update, err := s.today.GetToday(ctx, placeInfo.ID)
	if err != nil {
		if errors.Is(err, updates.ErrNotFound) {
			s.logger.Debug("no today content found for place", "place_id", placeInfo.ID)
			return &ContextResponse{Place: *placeInfo}, nil
		}
		return nil, fmt.Errorf("get today content: %w", err)
	}

	return &ContextResponse{
		Place:  *placeInfo,
		Update: update,
	}, nil
}

func (s *Service) resolvePlace(ctx context.Context, userID uuid.UUID, placeID *uuid.UUID) (*PlaceInfo, error) {
	if placeID != nil {
		p, err := s.places.ByID(ctx, *placeID)
		if err != nil {
			if errors.Is(err, places.ErrNotFound) {
				return nil, ErrPlaceNotFound
			}
			return nil, fmt.Errorf("lookup place: %w", err)
		}
		return &PlaceInfo{
			ID:        p.ID,
			Name:      p.Name,
			PlaceType: p.PlaceType,
			Lat:       p.Lat,
			Lon:       p.Lon,
		}, nil
	}

	saved, err := s.prefs.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list saved places: %w", err)
	}
	if len(saved) == 0 {
		return nil, ErrNoPlace
	}

	primary := saved[0]
	return &PlaceInfo{
		ID:        primary.ID,
		Name:      primary.Name,
		PlaceType: primary.PlaceType,
		Lat:       primary.Lat,
		Lon:       primary.Lon,
		Label:     primary.Label,
		IsPrimary: primary.IsPrimary,
	}, nil
}
