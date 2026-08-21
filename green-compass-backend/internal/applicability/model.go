package applicability

import (
	"time"

	"github.com/google/uuid"
)

type Rule struct {
	ID              uuid.UUID
	IndicatorID     uuid.UUID
	PlaceType       string
	Applicable      bool
	MinThreshold    *float64
	MaxThreshold    *float64
	RelevanceScore  float64 // 0.0–1.0
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type ApplicabilityResult struct {
	IndicatorID     uuid.UUID
	IndicatorCode   string
	DisplayName     string
	Applicable      bool
	RelevanceScore  float64
	MeetsThreshold  bool
	Value           float64
	Rule            *Rule
}

type FilterRequest struct {
	PlaceID     uuid.UUID
	PlaceType   string
	Indicators  []IndicatorData
}

type IndicatorData struct {
	IndicatorID uuid.UUID
	Code        string
	DisplayName string
	Value       float64
	Unit        string
}

type FilterResult struct {
	ApplicableIndicators []ApplicabilityResult
	FilteredCount        int
}
