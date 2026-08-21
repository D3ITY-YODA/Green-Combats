package normalization

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"green-compass-backend/internal/sources"
)

// RawInput is the subset of a landed raw record needed for normalization.
// Defined locally to avoid an import cycle with the ingestion package.
type RawInput struct {
	ID               uuid.UUID
	SourceID         uuid.UUID
	PlaceID          uuid.UUID
	Payload          []byte
	SourceObservedAt time.Time
	FetchedAt        time.Time
	SourceURL        string
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Normalize(ctx context.Context, sourceCode string, raw RawInput) ([]CanonicalObservation, error) {
	switch sourceCode {
	case sources.CodeOpenMeteo:
		return s.normalizeOpenMeteo(raw)
	case sources.CodeNASAPower:
		return s.normalizeNASAPower(raw)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedSource, sourceCode)
	}
}

func (s *Service) NormalizeAll(ctx context.Context, sourceCode string, records []RawInput) ([]CanonicalObservation, error) {
	var results []CanonicalObservation
	for _, rec := range records {
		obs, err := s.Normalize(ctx, sourceCode, rec)
		if err != nil {
			return nil, err
		}
		results = append(results, obs...)
	}
	return results, nil
}

func (s *Service) normalizeOpenMeteo(raw RawInput) ([]CanonicalObservation, error) {
	var payload map[string]json.RawMessage
	if err := json.Unmarshal(raw.Payload, &payload); err != nil {
		return nil, fmt.Errorf("%w: failed to decode Open-Meteo json: %v", ErrMalformedPayload, err)
	}

	mappings := []struct {
		key      string
		variable string
		unit     string
		min      float64
		max      float64
	}{
		{"temperature_2m", VarTemperature, UnitCelsius, -100, 70},
		{"relative_humidity_2m", VarRelativeHumidity, UnitPercent, 0, 100},
		{"precipitation", VarPrecipitation, UnitMillimeter, 0, 2000},
		{"wind_speed_10m", VarWindSpeed, UnitMetersPerSecond, 0, 150},
		{"soil_moisture_0_to_1cm", VarSoilMoisture, UnitM3PerM3, 0, 1.0},
		{"weather_code", VarWeatherCode, UnitWMOCode, 0, 100},
	}

	var results []CanonicalObservation
	for _, m := range mappings {
		rawVal, ok := payload[m.key]
		if !ok || len(rawVal) == 0 || string(rawVal) == "null" {
			continue
		}
		var val float64
		if err := json.Unmarshal(rawVal, &val); err != nil {
			continue
		}

		valid := val >= m.min && val <= m.max
		topicKey := inferTopicKey(m.variable)
		results = append(results, CanonicalObservation{
			ID:            uuid.New(),
			RawRecordID:   raw.ID,
			SourceID:      raw.SourceID,
			PlaceID:       raw.PlaceID,
			SourceCode:    sources.CodeOpenMeteo,
			DatasetKey:    "forecast",
			TopicKey:      topicKey,
			Variable:      m.variable,
			Value:         val,
			Unit:          m.unit,
			ObservedAt:    raw.SourceObservedAt,
			RetrievedAt:   raw.FetchedAt,
			IsForecast:    raw.SourceObservedAt.After(time.Now().UTC()),
			QualityStatus: "accepted",
			RawAssetRef:   raw.SourceURL,
			Valid:         valid,
		})
	}

	return results, nil
}

func (s *Service) normalizeNASAPower(raw RawInput) ([]CanonicalObservation, error) {
	var payload map[string]json.RawMessage
	if err := json.Unmarshal(raw.Payload, &payload); err != nil {
		return nil, fmt.Errorf("%w: failed to decode NASA POWER json: %v", ErrMalformedPayload, err)
	}

	mappings := []struct {
		key      string
		variable string
		unit     string
		min      float64
		max      float64
	}{
		{"T2M", VarTemperature, UnitCelsius, -100, 70},
		{"PRECTOTCORR", VarPrecipitation, UnitMillimeter, 0, 2000},
		{"RH2M", VarRelativeHumidity, UnitPercent, 0, 100},
		{"WS2M", VarWindSpeed, UnitMetersPerSecond, 0, 150},
		{"ALLSKY_SFC_SW_DWN", VarSolarRadiation, UnitMegajoulesPerM2, 0, 100},
	}

	var results []CanonicalObservation
	for _, m := range mappings {
		rawVal, ok := payload[m.key]
		if !ok || len(rawVal) == 0 || string(rawVal) == "null" {
			continue
		}
		var val float64
		if err := json.Unmarshal(rawVal, &val); err != nil {
			continue
		}

		// NASA POWER uses -999 as missing/null value sentinel
		if val == -999 {
			continue
		}

		valid := val >= m.min && val <= m.max
		topicKey := inferTopicKey(m.variable)
		results = append(results, CanonicalObservation{
			ID:            uuid.New(),
			RawRecordID:   raw.ID,
			SourceID:      raw.SourceID,
			PlaceID:       raw.PlaceID,
			SourceCode:    sources.CodeNASAPower,
			DatasetKey:    "historical",
			TopicKey:      topicKey,
			Variable:      m.variable,
			Value:         val,
			Unit:          m.unit,
			ObservedAt:    raw.SourceObservedAt,
			RetrievedAt:   raw.FetchedAt,
			IsForecast:    false,
			QualityStatus: "accepted",
			RawAssetRef:   raw.SourceURL,
			Valid:         valid,
		})
	}

	return results, nil
}
