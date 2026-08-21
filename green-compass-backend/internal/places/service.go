package places

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"green-compass-backend/pkg/geo"
)

const (
	maxNameLen         = 200
	defaultRadiusM     = 10_000
	maxRadiusM         = 50_000
	defaultNearbyLimit = 25
	maxNearbyLimit     = 100
)

type Caller struct {
	UserID          uuid.UUID
	IsPlatformAdmin bool
}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type CreateInput struct {
	Name      string
	PlaceType string
	Lat       float64
	Lon       float64
	Caller    Caller
}

func (s *Service) Create(ctx context.Context, in CreateInput) (*Place, error) {
	name, err := validateName(in.Name)
	if err != nil {
		return nil, err
	}
	if !IsValidType(in.PlaceType) {
		return nil, fmt.Errorf("%w: unknown place type %q", ErrInvalidData, in.PlaceType)
	}
	if _, err := geo.NewPoint(in.Lat, in.Lon); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidData, err)
	}
	if IsOfficialType(in.PlaceType) && !in.Caller.IsPlatformAdmin {
		return nil, fmt.Errorf("%w: only platform admins may create official places", ErrNotAllowed)
	}

	p := &Place{
		Name:      name,
		PlaceType: in.PlaceType,
		Lat:       in.Lat,
		Lon:       in.Lon,
		CreatedBy: &in.Caller.UserID,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) ByID(ctx context.Context, id uuid.UUID) (*Place, error) {
	return s.repo.ByID(ctx, id)
}

func (s *Service) Nearby(ctx context.Context, lat, lon, radiusM float64, limit int) ([]WithDistance, error) {
	if _, err := geo.NewPoint(lat, lon); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidData, err)
	}
	if radiusM == 0 {
		radiusM = defaultRadiusM
	}
	if limit == 0 {
		limit = defaultNearbyLimit
	}
	if radiusM < 0 || radiusM > maxRadiusM {
		return nil, fmt.Errorf("%w: radius must be between 1 and %d metres", ErrInvalidData, maxRadiusM)
	}
	if limit < 0 || limit > maxNearbyLimit {
		return nil, fmt.Errorf("%w: limit must be between 1 and %d", ErrInvalidData, maxNearbyLimit)
	}
	return s.repo.Nearby(ctx, lat, lon, radiusM, limit)
}

func (s *Service) Search(ctx context.Context, query string, limit int) ([]Place, error) {
	if limit <= 0 {
		limit = 25
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.Search(ctx, query, limit)
}

type UpdateInput struct {
	Name *string
	Lat  *float64
	Lon  *float64
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, in UpdateInput, caller Caller) (*Place, error) {
	if in.Name == nil && in.Lat == nil && in.Lon == nil {
		return nil, fmt.Errorf("%w: at least one field must be provided", ErrInvalidData)
	}
	existing, err := s.repo.ByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !canEdit(existing, caller) {
		return nil, ErrNotAllowed
	}

	if in.Name != nil {
		name, err := validateName(*in.Name)
		if err != nil {
			return nil, err
		}
		existing.Name = name
	}
	if in.Lat != nil || in.Lon != nil {
		if _, err := geo.NewPoint(*orDefault(in.Lat, existing.Lat), *orDefault(in.Lon, existing.Lon)); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidData, err)
		}
	}
	if in.Lat != nil {
		existing.Lat = *in.Lat
	}
	if in.Lon != nil {
		existing.Lon = *in.Lon
	}

	return s.repo.Update(ctx, existing)
}

func canEdit(p *Place, caller Caller) bool {
	if caller.IsPlatformAdmin {
		return true
	}
	if p.PlaceType != TypeCustom {
		return false
	}
	return p.CreatedBy != nil && *p.CreatedBy == caller.UserID
}

func validateName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", fmt.Errorf("%w: name is required", ErrInvalidData)
	}
	if len(name) > maxNameLen {
		return "", fmt.Errorf("%w: name exceeds %d characters", ErrInvalidData, maxNameLen)
	}
	return name, nil
}

func orDefault(v *float64, def float64) *float64 {
	if v != nil {
		return v
	}
	return &def
}
