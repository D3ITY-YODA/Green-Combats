package ingestion_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"green-compass-backend/internal/ingestion"
	"green-compass-backend/internal/places"
	"green-compass-backend/internal/sources"
	"green-compass-backend/internal/testutil"
)

func TestRepository_RawProvenanceAndDeduplication_Integration(t *testing.T) {
	pool := testutil.TestPool(t)
	ctx := context.Background()
	placeRepo := places.NewRepository(pool)
	place := &places.Place{Name: "Ingestion " + uuid.NewString(), PlaceType: places.TypeCommunity, Lat: -1.2, Lon: 36.8}
	if err := placeRepo.Create(ctx, place); err != nil {
		t.Fatalf("create place: %v", err)
	}
	source, err := sources.NewRepository(pool).ByCode(ctx, sources.CodeOpenMeteo)
	if err != nil {
		t.Fatalf("source: %v", err)
	}
	repo := ingestion.NewRepository(pool)
	run, err := repo.StartRun(ctx, source.ID, place.ID, nil, nil)
	if err != nil {
		t.Fatalf("StartRun: %v", err)
	}
	observedAt := time.Now().UTC().Truncate(time.Second)
	record := &ingestion.RawRecord{
		IngestionRunID: run.ID, SourceID: source.ID, PlaceID: place.ID, SourceObservedAt: observedAt,
		SourceURL:   "https://api.open-meteo.com/v1/forecast?latitude=-1.2&longitude=36.8",
		ContentType: "application/json", Payload: []byte(`{"hourly":{"temperature_2m":[21.5]}}`),
	}
	inserted, err := repo.StoreRaw(ctx, record)
	if err != nil {
		t.Fatalf("StoreRaw: %v", err)
	}
	if !inserted || record.ID == uuid.Nil || len(record.PayloadChecksum) != 64 {
		t.Errorf("stored record = %+v, inserted = %t", record, inserted)
	}

	duplicate := *record
	duplicate.ID = uuid.Nil
	inserted, err = repo.StoreRaw(ctx, &duplicate)
	if err != nil || inserted {
		t.Errorf("duplicate StoreRaw = inserted %t, err %v; want false, nil", inserted, err)
	}
	if err := repo.FinishRun(ctx, run.ID, ingestion.StatusSucceeded, 1, nil); err != nil {
		t.Fatalf("FinishRun: %v", err)
	}
	if err := repo.FinishRun(ctx, run.ID, ingestion.StatusSucceeded, 1, nil); !errors.Is(err, ingestion.ErrInvalidData) {
		t.Errorf("second FinishRun error = %v, want ErrInvalidData", err)
	}

	invalidRun, err := repo.StartRun(ctx, source.ID, place.ID, nil, nil)
	if err != nil {
		t.Fatalf("StartRun invalid: %v", err)
	}
	_, err = repo.StoreRaw(ctx, &ingestion.RawRecord{IngestionRunID: invalidRun.ID, SourceID: source.ID, PlaceID: place.ID, SourceObservedAt: observedAt, SourceURL: "https://example.test", ContentType: "application/json", Payload: []byte(`no`)})
	if !errors.Is(err, ingestion.ErrInvalidData) {
		t.Errorf("invalid JSON error = %v, want ErrInvalidData", err)
	}
}
