package indicators

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Definition struct {
	ID               uuid.UUID
	Code             string
	DisplayName      string
	Category         string // weather, water, agriculture, air_quality
	SourceVariables  pq.StringArray
	Description      string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type Indicator struct {
	ID                 uuid.UUID
	PlaceID            uuid.UUID
	IndicatorID        uuid.UUID
	ComputedAt         time.Time
	PeriodStart        time.Time
	PeriodEnd          time.Time
	Value              float64
	Unit               string
	Trend              *string // "increasing", "stable", "decreasing"
	TrendConfidence    *float64
	DataPointsCount    int
	Metadata           map[string]interface{}
	CreatedAt          time.Time
}

type IndicatorWithDefinition struct {
	*Indicator
	Definition *Definition
}

// ComputeRequest defines parameters for computing a single indicator
type ComputeRequest struct {
	PlaceID     uuid.UUID
	IndicatorID uuid.UUID
	PeriodStart time.Time
	PeriodEnd   time.Time
}

// ComputeResult holds the outcome of a computation
type ComputeResult struct {
	Indicator *Indicator
	Error     error
	Attempts  int
}
