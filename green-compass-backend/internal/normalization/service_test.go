package normalization_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"green-compass-backend/internal/normalization"
	"green-compass-backend/internal/sources"
)

func TestService_Normalize_OpenMeteo(t *testing.T) {
	svc := normalization.NewService()
	rawID := uuid.New()
	sourceID := uuid.New()
	placeID := uuid.New()
	observedAt := time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC)

	payload := map[string]any{
		"time":                   "2026-08-21T10:00",
		"temperature_2m":         24.5,
		"relative_humidity_2m":   68.0,
		"precipitation":          2.5,
		"wind_speed_10m":         5.2,
		"soil_moisture_0_to_1cm": 0.32,
		"weather_code":           61,
	}
	encoded, _ := json.Marshal(payload)

	raw := normalization.RawInput{
		ID:               rawID,
		SourceID:         sourceID,
		PlaceID:          placeID,
		SourceObservedAt: observedAt,
		Payload:          encoded,
	}

	obs, err := svc.Normalize(context.Background(), sources.CodeOpenMeteo, raw)
	if err != nil {
		t.Fatalf("Normalize failed: %v", err)
	}

	if len(obs) != 6 {
		t.Fatalf("expected 6 observations, got %d", len(obs))
	}

	byVar := make(map[string]normalization.CanonicalObservation)
	for _, o := range obs {
		byVar[o.Variable] = o
		if o.RawRecordID != rawID {
			t.Errorf("expected RawRecordID %v, got %v", rawID, o.RawRecordID)
		}
		if o.PlaceID != placeID {
			t.Errorf("expected PlaceID %v, got %v", placeID, o.PlaceID)
		}
		if !o.Valid {
			t.Errorf("expected observation %s to be valid", o.Variable)
		}
	}

	if temp, ok := byVar[normalization.VarTemperature]; !ok || temp.Value != 24.5 || temp.Unit != normalization.UnitCelsius {
		t.Errorf("unexpected temperature: %+v", temp)
	}
	if prec, ok := byVar[normalization.VarPrecipitation]; !ok || prec.Value != 2.5 || prec.Unit != normalization.UnitMillimeter {
		t.Errorf("unexpected precipitation: %+v", prec)
	}
	if rh, ok := byVar[normalization.VarRelativeHumidity]; !ok || rh.Value != 68.0 || rh.Unit != normalization.UnitPercent {
		t.Errorf("unexpected humidity: %+v", rh)
	}
	if ws, ok := byVar[normalization.VarWindSpeed]; !ok || ws.Value != 5.2 || ws.Unit != normalization.UnitMetersPerSecond {
		t.Errorf("unexpected wind speed: %+v", ws)
	}
	if sm, ok := byVar[normalization.VarSoilMoisture]; !ok || sm.Value != 0.32 || sm.Unit != normalization.UnitM3PerM3 {
		t.Errorf("unexpected soil moisture: %+v", sm)
	}
	if wc, ok := byVar[normalization.VarWeatherCode]; !ok || wc.Value != 61 || wc.Unit != normalization.UnitWMOCode {
		t.Errorf("unexpected weather code: %+v", wc)
	}
}

func TestService_Normalize_NASAPower(t *testing.T) {
	svc := normalization.NewService()
	rawID := uuid.New()
	sourceID := uuid.New()
	placeID := uuid.New()
	observedAt := time.Date(2026, 8, 21, 0, 0, 0, 0, time.UTC)

	payload := map[string]any{
		"date":              "20260821",
		"T2M":               22.8,
		"PRECTOTCORR":       0.0,
		"RH2M":              75.5,
		"WS2M":              3.1,
		"ALLSKY_SFC_SW_DWN": 19.4,
	}
	encoded, _ := json.Marshal(payload)

	raw := normalization.RawInput{
		ID:               rawID,
		SourceID:         sourceID,
		PlaceID:          placeID,
		SourceObservedAt: observedAt,
		Payload:          encoded,
	}

	obs, err := svc.Normalize(context.Background(), sources.CodeNASAPower, raw)
	if err != nil {
		t.Fatalf("Normalize failed: %v", err)
	}

	if len(obs) != 5 {
		t.Fatalf("expected 5 observations, got %d", len(obs))
	}

	byVar := make(map[string]normalization.CanonicalObservation)
	for _, o := range obs {
		byVar[o.Variable] = o
	}

	if rad, ok := byVar[normalization.VarSolarRadiation]; !ok || rad.Value != 19.4 || rad.Unit != normalization.UnitMegajoulesPerM2 {
		t.Errorf("unexpected solar radiation: %+v", rad)
	}
}

func TestService_Normalize_NASAPower_MissingSentinel(t *testing.T) {
	svc := normalization.NewService()
	payload := map[string]any{
		"date":              "20260821",
		"T2M":               20.0,
		"PRECTOTCORR":       -999.0, // sentinel missing value in NASA POWER
		"RH2M":              70.0,
		"WS2M":              2.0,
		"ALLSKY_SFC_SW_DWN": -999.0,
	}
	encoded, _ := json.Marshal(payload)

	raw := normalization.RawInput{Payload: encoded}
	obs, err := svc.Normalize(context.Background(), sources.CodeNASAPower, raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should omit -999 missing entries (leaving 3 valid variables)
	if len(obs) != 3 {
		t.Fatalf("expected 3 valid observations, got %d", len(obs))
	}
}

func TestService_Normalize_Errors(t *testing.T) {
	svc := normalization.NewService()

	t.Run("unsupported source", func(t *testing.T) {
		raw := normalization.RawInput{Payload: json.RawMessage(`{}`)}
		_, err := svc.Normalize(context.Background(), "unknown_source", raw)
		if !errors.Is(err, normalization.ErrUnsupportedSource) {
			t.Errorf("expected ErrUnsupportedSource, got %v", err)
		}
	})

	t.Run("malformed json payload", func(t *testing.T) {
		raw := normalization.RawInput{Payload: json.RawMessage(`invalid-json`)}
		_, err := svc.Normalize(context.Background(), sources.CodeOpenMeteo, raw)
		if !errors.Is(err, normalization.ErrMalformedPayload) {
			t.Errorf("expected ErrMalformedPayload, got %v", err)
		}
	})
}
