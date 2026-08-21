package ingestion

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"green-compass-backend/internal/connectors"
	"green-compass-backend/internal/places"
	"green-compass-backend/internal/sources"
)

var (
	ErrConnectorNotFound = errors.New("connector not found for source")
	ErrSourceDisabled    = errors.New("data source is disabled")
	ErrInvalidRequest    = errors.New("invalid ingestion request")
)

type RunRepository interface {
	StartRun(ctx context.Context, sourceID, placeID uuid.UUID, from, to *time.Time) (*Run, error)
	StoreRaw(ctx context.Context, record *RawRecord) (inserted bool, err error)
	FinishRun(ctx context.Context, id uuid.UUID, status string, landedCount int, message *string) error
}

type Service struct {
	repo     RunRepository
	registry *connectors.Registry
}

func NewService(repo RunRepository, registry *connectors.Registry) *Service {
	return &Service{
		repo:     repo,
		registry: registry,
	}
}

type IngestRequest struct {
	Source sources.Source
	Place  places.Place
	From   time.Time
	To     time.Time
}

type RunResult struct {
	RunID        uuid.UUID
	SourceCode   string
	PlaceID      uuid.UUID
	FetchedCount int
	LandedCount  int
	Status       string
	ErrorMessage *string
}

func (s *Service) Ingest(ctx context.Context, req IngestRequest) (*RunResult, error) {
	if !req.Source.Enabled {
		return nil, fmt.Errorf("%w: %s", ErrSourceDisabled, req.Source.Code)
	}
	if req.Place.ID == uuid.Nil {
		return nil, fmt.Errorf("%w: place ID is required", ErrInvalidRequest)
	}
	if req.From.IsZero() || req.To.IsZero() || req.To.Before(req.From) {
		return nil, fmt.Errorf("%w: invalid time window", ErrInvalidRequest)
	}

	connector, ok := s.registry.ByCode(req.Source.Code)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrConnectorNotFound, req.Source.Code)
	}

	run, err := s.repo.StartRun(ctx, req.Source.ID, req.Place.ID, &req.From, &req.To)
	if err != nil {
		return nil, fmt.Errorf("start ingestion run: %w", err)
	}

	observations, err := connector.Fetch(ctx, connectors.FetchRequest{
		PlaceID: req.Place.ID,
		Lat:     req.Place.Lat,
		Lon:     req.Place.Lon,
		From:    req.From,
		To:      req.To,
	})
	if err != nil {
		errMsg := err.Error()
		_ = s.repo.FinishRun(ctx, run.ID, StatusFailed, 0, &errMsg)
		return &RunResult{
			RunID:        run.ID,
			SourceCode:   req.Source.Code,
			PlaceID:      req.Place.ID,
			FetchedCount: 0,
			LandedCount:  0,
			Status:       StatusFailed,
			ErrorMessage: &errMsg,
		}, fmt.Errorf("fetch observations from connector %s: %w", req.Source.Code, err)
	}

	landedCount := 0
	for _, obs := range observations {
		record := &RawRecord{
			IngestionRunID:   run.ID,
			SourceID:         req.Source.ID,
			PlaceID:          req.Place.ID,
			SourceObservedAt: obs.ObservedAt,
			SourceURL:        obs.SourceURL,
			ContentType:      obs.ContentType,
			Payload:          obs.Payload,
			PayloadChecksum:  Checksum(obs.Payload),
		}

		inserted, err := s.repo.StoreRaw(ctx, record)
		if err != nil {
			errMsg := fmt.Sprintf("store raw observation: %v", err)
			_ = s.repo.FinishRun(ctx, run.ID, StatusFailed, landedCount, &errMsg)
			return &RunResult{
				RunID:        run.ID,
				SourceCode:   req.Source.Code,
				PlaceID:      req.Place.ID,
				FetchedCount: len(observations),
				LandedCount:  landedCount,
				Status:       StatusFailed,
				ErrorMessage: &errMsg,
			}, err
		}
		if inserted {
			landedCount++
		}
	}

	if err := s.repo.FinishRun(ctx, run.ID, StatusSucceeded, landedCount, nil); err != nil {
		return nil, fmt.Errorf("finish ingestion run: %w", err)
	}

	return &RunResult{
		RunID:        run.ID,
		SourceCode:   req.Source.Code,
		PlaceID:      req.Place.ID,
		FetchedCount: len(observations),
		LandedCount:  landedCount,
		Status:       StatusSucceeded,
	}, nil
}
