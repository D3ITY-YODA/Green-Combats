package ingestion_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"green-compass-backend/internal/connectors"
	"green-compass-backend/internal/ingestion"
	"green-compass-backend/internal/places"
	"green-compass-backend/internal/sources"
)

type mockRunRepo struct {
	startRunFunc  func(ctx context.Context, sourceID, placeID uuid.UUID, from, to *time.Time) (*ingestion.Run, error)
	storeRawFunc  func(ctx context.Context, record *ingestion.RawRecord) (bool, error)
	finishRunFunc func(ctx context.Context, id uuid.UUID, status string, landedCount int, message *string) error
}

func (m *mockRunRepo) StartRun(ctx context.Context, sourceID, placeID uuid.UUID, from, to *time.Time) (*ingestion.Run, error) {
	if m.startRunFunc != nil {
		return m.startRunFunc(ctx, sourceID, placeID, from, to)
	}
	return &ingestion.Run{ID: uuid.New(), SourceID: sourceID, PlaceID: placeID, Status: ingestion.StatusRunning}, nil
}

func (m *mockRunRepo) StoreRaw(ctx context.Context, record *ingestion.RawRecord) (bool, error) {
	if m.storeRawFunc != nil {
		return m.storeRawFunc(ctx, record)
	}
	return true, nil
}

func (m *mockRunRepo) FinishRun(ctx context.Context, id uuid.UUID, status string, landedCount int, message *string) error {
	if m.finishRunFunc != nil {
		return m.finishRunFunc(ctx, id, status, landedCount, message)
	}
	return nil
}

type mockConnector struct {
	code      string
	fetchFunc func(ctx context.Context, request connectors.FetchRequest) ([]connectors.RawObservation, error)
}

func (m *mockConnector) Code() string { return m.code }
func (m *mockConnector) Fetch(ctx context.Context, request connectors.FetchRequest) ([]connectors.RawObservation, error) {
	if m.fetchFunc != nil {
		return m.fetchFunc(ctx, request)
	}
	return nil, nil
}

func TestService_Ingest_Validation(t *testing.T) {
	placeID := uuid.New()
	validPlace := places.Place{ID: placeID, Lat: 1.23, Lon: 36.8}
	validSource := sources.Source{ID: uuid.New(), Code: "open_meteo", Enabled: true}
	now := time.Now().UTC()

	tests := []struct {
		name      string
		req       ingestion.IngestRequest
		wantErrIs error
		connector connectors.Connector
	}{
		{
			name: "disabled source",
			req: ingestion.IngestRequest{
				Source: sources.Source{ID: uuid.New(), Code: "open_meteo", Enabled: false},
				Place:  validPlace,
				From:   now.Add(-24 * time.Hour),
				To:     now,
			},
			wantErrIs: ingestion.ErrSourceDisabled,
		},
		{
			name: "nil place ID",
			req: ingestion.IngestRequest{
				Source: validSource,
				Place:  places.Place{},
				From:   now.Add(-24 * time.Hour),
				To:     now,
			},
			wantErrIs: ingestion.ErrInvalidRequest,
		},
		{
			name: "zero from time",
			req: ingestion.IngestRequest{
				Source: validSource,
				Place:  validPlace,
				From:   time.Time{},
				To:     now,
			},
			wantErrIs: ingestion.ErrInvalidRequest,
		},
		{
			name: "to before from",
			req: ingestion.IngestRequest{
				Source: validSource,
				Place:  validPlace,
				From:   now,
				To:     now.Add(-time.Hour),
			},
			wantErrIs: ingestion.ErrInvalidRequest,
		},
		{
			name: "unregistered connector",
			req: ingestion.IngestRequest{
				Source: sources.Source{ID: uuid.New(), Code: "unknown_source", Enabled: true},
				Place:  validPlace,
				From:   now.Add(-24 * time.Hour),
				To:     now,
			},
			wantErrIs: ingestion.ErrConnectorNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var conn connectors.Connector = &mockConnector{code: "open_meteo"}
			if tt.connector != nil {
				conn = tt.connector
			}
			reg, err := connectors.NewRegistry(conn)
			if err != nil {
				t.Fatalf("NewRegistry failed: %v", err)
			}
			repo := &mockRunRepo{}
			svc := ingestion.NewService(repo, reg, nil, nil)

			_, err = svc.Ingest(context.Background(), tt.req)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !errors.Is(err, tt.wantErrIs) {
				t.Fatalf("expected error is %v, got %v", tt.wantErrIs, err)
			}
		})
	}
}

func TestService_Ingest_Success(t *testing.T) {
	placeID := uuid.New()
	sourceID := uuid.New()
	place := places.Place{ID: placeID, Lat: -1.286389, Lon: 36.817223}
	source := sources.Source{ID: sourceID, Code: "open_meteo", Enabled: true}
	now := time.Now().UTC()
	from := now.Add(-24 * time.Hour)
	to := now

	obsData := []connectors.RawObservation{
		{ObservedAt: from.Add(time.Hour), SourceURL: "https://api.test/1", ContentType: "application/json", Payload: json.RawMessage(`{"temp":25.5}`)},
		{ObservedAt: from.Add(2 * time.Hour), SourceURL: "https://api.test/2", ContentType: "application/json", Payload: json.RawMessage(`{"temp":26.0}`)},
		{ObservedAt: from.Add(3 * time.Hour), SourceURL: "https://api.test/3", ContentType: "application/json", Payload: json.RawMessage(`{"temp":24.8}`)},
	}

	conn := &mockConnector{
		code: "open_meteo",
		fetchFunc: func(ctx context.Context, request connectors.FetchRequest) ([]connectors.RawObservation, error) {
			if request.PlaceID != placeID {
				t.Errorf("expected placeID %v, got %v", placeID, request.PlaceID)
			}
			return obsData, nil
		},
	}
	reg, _ := connectors.NewRegistry(conn)

	var storedRecords []*ingestion.RawRecord
	var finishStatus string
	var finishCount int

	repo := &mockRunRepo{
		storeRawFunc: func(ctx context.Context, record *ingestion.RawRecord) (bool, error) {
			storedRecords = append(storedRecords, record)
			return true, nil
		},
		finishRunFunc: func(ctx context.Context, id uuid.UUID, status string, landedCount int, message *string) error {
			finishStatus = status
			finishCount = landedCount
			return nil
		},
	}

	svc := ingestion.NewService(repo, reg, nil, nil)
	result, err := svc.Ingest(context.Background(), ingestion.IngestRequest{
		Source: source,
		Place:  place,
		From:   from,
		To:     to,
	})

	if err != nil {
		t.Fatalf("Ingest failed: %v", err)
	}
	if result.Status != ingestion.StatusSucceeded {
		t.Errorf("expected status %s, got %s", ingestion.StatusSucceeded, result.Status)
	}
	if result.FetchedCount != 3 {
		t.Errorf("expected 3 fetched, got %d", result.FetchedCount)
	}
	if result.LandedCount != 3 {
		t.Errorf("expected 3 landed, got %d", result.LandedCount)
	}
	if len(storedRecords) != 3 {
		t.Errorf("expected 3 stored records, got %d", len(storedRecords))
	}
	if finishStatus != ingestion.StatusSucceeded {
		t.Errorf("expected finish status succeeded, got %s", finishStatus)
	}
	if finishCount != 3 {
		t.Errorf("expected finish count 3, got %d", finishCount)
	}
}

func TestService_Ingest_Deduplication(t *testing.T) {
	placeID := uuid.New()
	sourceID := uuid.New()
	place := places.Place{ID: placeID, Lat: 0, Lon: 37}
	source := sources.Source{ID: sourceID, Code: "open_meteo", Enabled: true}
	now := time.Now().UTC()

	conn := &mockConnector{
		code: "open_meteo",
		fetchFunc: func(ctx context.Context, request connectors.FetchRequest) ([]connectors.RawObservation, error) {
			return []connectors.RawObservation{
				{ObservedAt: now, Payload: json.RawMessage(`{"val":1}`)},
				{ObservedAt: now.Add(time.Hour), Payload: json.RawMessage(`{"val":2}`)},
			}, nil
		},
	}
	reg, _ := connectors.NewRegistry(conn)

	callCount := 0
	repo := &mockRunRepo{
		storeRawFunc: func(ctx context.Context, record *ingestion.RawRecord) (bool, error) {
			callCount++
			if callCount == 2 {
				// Second record is a duplicate
				return false, nil
			}
			return true, nil
		},
	}

	svc := ingestion.NewService(repo, reg, nil, nil)
	result, err := svc.Ingest(context.Background(), ingestion.IngestRequest{
		Source: source,
		Place:  place,
		From:   now.Add(-2 * time.Hour),
		To:     now,
	})

	if err != nil {
		t.Fatalf("Ingest failed: %v", err)
	}
	if result.FetchedCount != 2 {
		t.Errorf("expected 2 fetched, got %d", result.FetchedCount)
	}
	if result.LandedCount != 1 {
		t.Errorf("expected 1 landed, got %d", result.LandedCount)
	}
}

func TestService_Ingest_ConnectorFetchError(t *testing.T) {
	placeID := uuid.New()
	sourceID := uuid.New()
	place := places.Place{ID: placeID, Lat: 0, Lon: 37}
	source := sources.Source{ID: sourceID, Code: "open_meteo", Enabled: true}
	now := time.Now().UTC()

	conn := &mockConnector{
		code: "open_meteo",
		fetchFunc: func(ctx context.Context, request connectors.FetchRequest) ([]connectors.RawObservation, error) {
			return nil, errors.New("upstream timeout")
		},
	}
	reg, _ := connectors.NewRegistry(conn)

	var finishedWithStatus string
	var finishedMessage *string

	repo := &mockRunRepo{
		finishRunFunc: func(ctx context.Context, id uuid.UUID, status string, landedCount int, message *string) error {
			finishedWithStatus = status
			finishedMessage = message
			return nil
		},
	}

	svc := ingestion.NewService(repo, reg, nil, nil)
	result, err := svc.Ingest(context.Background(), ingestion.IngestRequest{
		Source: source,
		Place:  place,
		From:   now.Add(-time.Hour),
		To:     now,
	})

	if err == nil {
		t.Fatalf("expected error from upstream, got nil")
	}
	if result == nil || result.Status != ingestion.StatusFailed {
		t.Errorf("expected result status failed, got %v", result)
	}
	if finishedWithStatus != ingestion.StatusFailed {
		t.Errorf("expected finish status failed, got %s", finishedWithStatus)
	}
	if finishedMessage == nil || *finishedMessage != "upstream timeout" {
		t.Errorf("expected error message 'upstream timeout', got %v", finishedMessage)
	}
}
