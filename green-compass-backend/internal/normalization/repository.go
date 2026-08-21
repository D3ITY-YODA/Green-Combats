package normalization

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Store(ctx context.Context, obs []CanonicalObservation) error
	ListByPlace(ctx context.Context, placeID uuid.UUID, from, to time.Time) ([]CanonicalObservation, error)
}
