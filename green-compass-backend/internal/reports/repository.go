package reports

import (
	"context"
	"time"

	"github.com/google/uuid"

	"green-compass-backend/internal/observations"
)

type Repository struct {
	obsRepo *observations.Repository
}

func NewRepository(obsRepo *observations.Repository) *Repository {
	return &Repository{obsRepo: obsRepo}
}

func (r *Repository) ListPending(ctx context.Context, limit int) ([]Report, error) {
	obs, err := r.obsRepo.ListPending(ctx, limit)
	if err != nil {
		return nil, err
	}
	reports := make([]Report, len(obs))
	for i, o := range obs {
		reports[i] = toReport(o)
	}
	return reports, nil
}

func (r *Repository) ByID(ctx context.Context, id uuid.UUID) (*Report, error) {
	o, err := r.obsRepo.ByID(ctx, id)
	if err != nil {
		return nil, err
	}
	report := toReport(*o)
	return &report, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, verifiedBy uuid.UUID) (*Report, error) {
	now := time.Now().UTC()
	o, err := r.obsRepo.UpdateStatus(ctx, id, status, &verifiedBy, &now)
	if err != nil {
		return nil, err
	}
	report := toReport(*o)
	return &report, nil
}

func toReport(o observations.Observation) Report {
	return Report{
		ID:             o.ID,
		ReporterID:     o.ReporterID,
		PlaceID:        o.PlaceID,
		Lat:            o.Lat,
		Lon:            o.Lon,
		Category:       o.Category,
		Description:    o.Description,
		PhotoObjectKey: o.PhotoObjectKey,
		Status:         o.Status,
		VerifiedBy:     o.VerifiedBy,
		VerifiedAt:     o.VerifiedAt,
		CreatedAt:      o.CreatedAt,
		UpdatedAt:      o.UpdatedAt,
	}
}
