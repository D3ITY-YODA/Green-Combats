package indicators

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"green-compass-backend/internal/normalization"
)

type mockNormRepo struct {
	observations []normalization.CanonicalObservation
}

func (m *mockNormRepo) ListCanonicalObservations(ctx context.Context, placeID, sourceID uuid.UUID, from, to time.Time) ([]normalization.CanonicalObservation, error) {
	return m.observations, nil
}

func newMockNormRepo(obs []normalization.CanonicalObservation) *mockNormRepo {
	return &mockNormRepo{observations: obs}
}

func TestRainIntensityTrend(t *testing.T) {
	tests := []struct {
		name       string
		obs        []normalization.CanonicalObservation
		wantValue  float64
		wantTrend  string
		wantPoints int
		delta      float64
	}{
		{
			name: "increasing rain trend",
			obs: []normalization.CanonicalObservation{
				{Variable: normalization.VarPrecipitation, Value: 1.0, Valid: true, ObservedAt: time.Now().Add(-7 * 24 * time.Hour)},
				{Variable: normalization.VarPrecipitation, Value: 1.5, Valid: true, ObservedAt: time.Now().Add(-6 * 24 * time.Hour)},
				{Variable: normalization.VarPrecipitation, Value: 2.0, Valid: true, ObservedAt: time.Now().Add(-5 * 24 * time.Hour)},
				{Variable: normalization.VarPrecipitation, Value: 5.0, Valid: true, ObservedAt: time.Now().Add(-2 * 24 * time.Hour)},
				{Variable: normalization.VarPrecipitation, Value: 6.0, Valid: true, ObservedAt: time.Now().Add(-1 * 24 * time.Hour)},
				{Variable: normalization.VarPrecipitation, Value: 7.0, Valid: true, ObservedAt: time.Now()},
			},
			wantValue:  6.0,
			wantTrend:  "increasing",
			wantPoints: 6,
			delta:      0.01,
		},
		{
			name: "decreasing rain trend",
			obs: []normalization.CanonicalObservation{
				{Variable: normalization.VarPrecipitation, Value: 10.0, Valid: true, ObservedAt: time.Now().Add(-5 * 24 * time.Hour)},
				{Variable: normalization.VarPrecipitation, Value: 8.0, Valid: true, ObservedAt: time.Now().Add(-4 * 24 * time.Hour)},
				{Variable: normalization.VarPrecipitation, Value: 1.0, Valid: true, ObservedAt: time.Now().Add(-3 * 24 * time.Hour)},
				{Variable: normalization.VarPrecipitation, Value: 0.5, Valid: true, ObservedAt: time.Now().Add(-2 * 24 * time.Hour)},
				{Variable: normalization.VarPrecipitation, Value: 0.2, Valid: true, ObservedAt: time.Now().Add(-1 * 24 * time.Hour)},
				{Variable: normalization.VarPrecipitation, Value: 0.1, Valid: true, ObservedAt: time.Now()},
			},
			wantValue:  0.27,
			wantTrend:  "decreasing",
			wantPoints: 6,
			delta:      0.01,
		},
		{
			name: "stable rain trend",
			obs: []normalization.CanonicalObservation{
				{Variable: normalization.VarPrecipitation, Value: 2.0, Valid: true, ObservedAt: time.Now().Add(-3 * 24 * time.Hour)},
				{Variable: normalization.VarPrecipitation, Value: 2.1, Valid: true, ObservedAt: time.Now().Add(-2 * 24 * time.Hour)},
				{Variable: normalization.VarPrecipitation, Value: 1.9, Valid: true, ObservedAt: time.Now().Add(-1 * 24 * time.Hour)},
				{Variable: normalization.VarPrecipitation, Value: 2.0, Valid: true, ObservedAt: time.Now()},
			},
			wantValue:  1.95,
			wantTrend:  "stable",
			wantPoints: 4,
			delta:      0.01,
		},
		{
			name:       "no data",
			obs:        []normalization.CanonicalObservation{},
			wantValue:  0,
			wantTrend:  "",
			wantPoints: 0,
			delta:      0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ind := &Indicator{
				ID:       uuid.New(),
				Metadata: make(map[string]interface{}),
			}
			svc := &Service{logger: slog.Default()}
			svc.computeRainIntensityTrend(tt.obs, ind)

			if ind.DataPointsCount != tt.wantPoints {
				t.Errorf("DataPointsCount: got %d, want %d", ind.DataPointsCount, tt.wantPoints)
			}
			if tt.wantPoints > 0 && delta(ind.Value, tt.wantValue) > tt.delta {
				t.Errorf("Value: got %v, want %v (delta %v)", ind.Value, tt.wantValue, delta(ind.Value, tt.wantValue))
			}
			if tt.wantPoints > 1 && tt.wantTrend != "" {
				if ind.Trend == nil {
					t.Errorf("Trend: got nil, want %s", tt.wantTrend)
				} else if *ind.Trend != tt.wantTrend {
					t.Errorf("Trend: got %s, want %s", *ind.Trend, tt.wantTrend)
				}
			}
		})
	}
}

func TestRainFrequency(t *testing.T) {
	tests := []struct {
		name       string
		obs        []normalization.CanonicalObservation
		wantValue  float64
		wantPoints int
	}{
		{
			name: "4 rainy days out of 6",
			obs: []normalization.CanonicalObservation{
				{Variable: normalization.VarPrecipitation, Value: 5.0, Valid: true},
				{Variable: normalization.VarPrecipitation, Value: 0.2, Valid: true},
				{Variable: normalization.VarPrecipitation, Value: 10.0, Valid: true},
				{Variable: normalization.VarPrecipitation, Value: 0.1, Valid: true},
				{Variable: normalization.VarPrecipitation, Value: 3.0, Valid: true},
				{Variable: normalization.VarPrecipitation, Value: 0.0, Valid: true},
			},
			wantValue:  3, // 5.0, 10.0, 3.0 are > 0.5 threshold
			wantPoints: 6,
		},
		{
			name: "no rainy days",
			obs: []normalization.CanonicalObservation{
				{Variable: normalization.VarPrecipitation, Value: 0.0, Valid: true},
				{Variable: normalization.VarPrecipitation, Value: 0.1, Valid: true},
				{Variable: normalization.VarPrecipitation, Value: 0.2, Valid: true},
			},
			wantValue:  0,
			wantPoints: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ind := &Indicator{
				ID:       uuid.New(),
				Metadata: make(map[string]interface{}),
			}
			svc := &Service{logger: slog.Default()}
			svc.computeRainFrequency(tt.obs, ind)

			if ind.Value != tt.wantValue {
				t.Errorf("Value: got %v, want %v", ind.Value, tt.wantValue)
			}
			if ind.DataPointsCount != tt.wantPoints {
				t.Errorf("DataPointsCount: got %d, want %d", ind.DataPointsCount, tt.wantPoints)
			}
		})
	}
}

func TestTemperatureTrend(t *testing.T) {
	tests := []struct {
		name       string
		obs        []normalization.CanonicalObservation
		wantTrend  *string
		wantPoints int
	}{
		{
			name: "warming trend",
			obs: []normalization.CanonicalObservation{
				{Variable: normalization.VarTemperature, Value: 15.0, Valid: true, ObservedAt: time.Now().Add(-4 * time.Hour)},
				{Variable: normalization.VarTemperature, Value: 16.0, Valid: true, ObservedAt: time.Now().Add(-3 * time.Hour)},
				{Variable: normalization.VarTemperature, Value: 25.0, Valid: true, ObservedAt: time.Now().Add(-2 * time.Hour)},
				{Variable: normalization.VarTemperature, Value: 26.0, Valid: true, ObservedAt: time.Now().Add(-1 * time.Hour)},
			},
			wantTrend:  strPtr("increasing"),
			wantPoints: 4,
		},
		{
			name: "cooling trend",
			obs: []normalization.CanonicalObservation{
				{Variable: normalization.VarTemperature, Value: 30.0, Valid: true, ObservedAt: time.Now().Add(-4 * time.Hour)},
				{Variable: normalization.VarTemperature, Value: 28.0, Valid: true, ObservedAt: time.Now().Add(-3 * time.Hour)},
				{Variable: normalization.VarTemperature, Value: 15.0, Valid: true, ObservedAt: time.Now().Add(-2 * time.Hour)},
				{Variable: normalization.VarTemperature, Value: 14.0, Valid: true, ObservedAt: time.Now().Add(-1 * time.Hour)},
			},
			wantTrend:  strPtr("decreasing"),
			wantPoints: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ind := &Indicator{
				ID:       uuid.New(),
				Metadata: make(map[string]interface{}),
			}
			svc := &Service{logger: slog.Default()}
			svc.computeTemperatureTrend(tt.obs, ind)

			if ind.DataPointsCount != tt.wantPoints {
				t.Errorf("DataPointsCount: got %d, want %d", ind.DataPointsCount, tt.wantPoints)
			}
			if (tt.wantTrend == nil) != (ind.Trend == nil) {
				t.Errorf("Trend: got %v, want %v", ind.Trend, tt.wantTrend)
			}
			if tt.wantTrend != nil && ind.Trend != nil && *ind.Trend != *tt.wantTrend {
				t.Errorf("Trend: got %s, want %s", *ind.Trend, *tt.wantTrend)
			}
		})
	}
}

func TestSoilMoistureLevel(t *testing.T) {
	tests := []struct {
		name       string
		obs        []normalization.CanonicalObservation
		wantValue  float64
		delta      float64
		wantPoints int
	}{
		{
			name: "40% soil moisture",
			obs: []normalization.CanonicalObservation{
				{Variable: normalization.VarSoilMoisture, Value: 0.4, Valid: true},
				{Variable: normalization.VarSoilMoisture, Value: 0.4, Valid: true},
				{Variable: normalization.VarSoilMoisture, Value: 0.4, Valid: true},
			},
			wantValue:  40.0,
			delta:      0.1,
			wantPoints: 3,
		},
		{
			name: "70% soil moisture",
			obs: []normalization.CanonicalObservation{
				{Variable: normalization.VarSoilMoisture, Value: 0.7, Valid: true},
				{Variable: normalization.VarSoilMoisture, Value: 0.7, Valid: true},
			},
			wantValue:  70.0,
			delta:      0.1,
			wantPoints: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ind := &Indicator{
				ID:       uuid.New(),
				Metadata: make(map[string]interface{}),
			}
			svc := &Service{logger: slog.Default()}
			svc.computeSoilMoistureLevel(tt.obs, ind)

			if ind.DataPointsCount != tt.wantPoints {
				t.Errorf("DataPointsCount: got %d, want %d", ind.DataPointsCount, tt.wantPoints)
			}
			if delta(ind.Value, tt.wantValue) > tt.delta {
				t.Errorf("Value: got %v, want %v (delta %v)", ind.Value, tt.wantValue, delta(ind.Value, tt.wantValue))
			}
		})
	}
}

func TestDroughtStress(t *testing.T) {
	tests := []struct {
		name       string
		obs        []normalization.CanonicalObservation
		wantMin    float64
		wantMax    float64
		wantPoints int
	}{
		{
			name: "normal conditions (low stress)",
			obs: []normalization.CanonicalObservation{
				{Variable: normalization.VarPrecipitation, Value: 40.0, Valid: true},
				{Variable: normalization.VarTemperature, Value: 20.0, Valid: true},
				{Variable: normalization.VarSoilMoisture, Value: 0.5, Valid: true},
			},
			wantMin:    0,
			wantMax:    30,
			wantPoints: 1,
		},
		{
			name: "drought conditions (high stress)",
			obs: []normalization.CanonicalObservation{
				{Variable: normalization.VarPrecipitation, Value: 5.0, Valid: true},
				{Variable: normalization.VarTemperature, Value: 35.0, Valid: true},
				{Variable: normalization.VarSoilMoisture, Value: 0.1, Valid: true},
			},
			wantMin:    70,
			wantMax:    100,
			wantPoints: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ind := &Indicator{
				ID:       uuid.New(),
				Metadata: make(map[string]interface{}),
			}
			svc := &Service{logger: slog.Default()}
			svc.computeDroughtStress(tt.obs, ind)

			if ind.Value < tt.wantMin || ind.Value > tt.wantMax {
				t.Errorf("Value: got %v, want in range [%v, %v]", ind.Value, tt.wantMin, tt.wantMax)
			}
		})
	}
}

// Helpers
func strPtr(s string) *string {
	return &s
}

func delta(a, b float64) float64 {
	if a > b {
		return a - b
	}
	return b - a
}
