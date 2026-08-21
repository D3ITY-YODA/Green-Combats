package preferences

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service { return &Service{repo: repo} }

func (s *Service) List(ctx context.Context, userID uuid.UUID) ([]SavedPlace, error) {
	return s.repo.List(ctx, userID)
}

func (s *Service) Save(ctx context.Context, userID, placeID uuid.UUID, label *string) error {
	label, err := normalizeLabel(label)
	if err != nil {
		return err
	}
	return s.repo.Save(ctx, userID, placeID, label)
}

func (s *Service) Unsave(ctx context.Context, userID, placeID uuid.UUID) error {
	return s.repo.Unsave(ctx, userID, placeID)
}

func (s *Service) SetPrimary(ctx context.Context, userID, placeID uuid.UUID) error {
	return s.repo.SetPrimary(ctx, userID, placeID)
}

func normalizeLabel(label *string) (*string, error) {
	if label == nil {
		return nil, nil
	}
	value := strings.TrimSpace(*label)
	if value == "" {
		return nil, nil
	}
	if len(value) > MaxLabelLen {
		return nil, fmt.Errorf("%w: label exceeds %d characters", ErrInvalidData, MaxLabelLen)
	}
	return &value, nil
}
