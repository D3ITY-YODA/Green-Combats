package sources_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"green-compass-backend/internal/sources"
	"green-compass-backend/internal/testutil"
)

func TestRepository_SeededSources_Integration(t *testing.T) {
	repo := sources.NewRepository(testutil.TestPool(t))
	ctx := context.Background()

	tests := []struct {
		code     string
		interval time.Duration
		enabled  bool
	}{
		{sources.CodeOpenMeteo, 6 * time.Hour, true},
		{sources.CodeNASAPower, 24 * time.Hour, true},
		{sources.CodeFreshwater, 6 * time.Hour, false},
		{sources.CodeSatellite, 24 * time.Hour, false},
		{sources.CodeCommunity, time.Hour, false},
	}
	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			source, err := repo.ByCode(ctx, tt.code)
			if err != nil {
				t.Fatalf("ByCode: %v", err)
			}
			if source.PollInterval != tt.interval || source.Enabled != tt.enabled {
				t.Errorf("source = %+v, want interval %s enabled %t", source, tt.interval, tt.enabled)
			}
		})
	}
	if _, err := repo.ByCode(ctx, "missing"); !errors.Is(err, sources.ErrNotFound) {
		t.Errorf("unknown source error = %v, want ErrNotFound", err)
	}

	enabled, err := repo.Enabled(ctx)
	if err != nil {
		t.Fatalf("Enabled: %v", err)
	}
	if len(enabled) != 2 || enabled[0].Code != sources.CodeNASAPower || enabled[1].Code != sources.CodeOpenMeteo {
		t.Errorf("enabled sources = %+v", enabled)
	}
}
