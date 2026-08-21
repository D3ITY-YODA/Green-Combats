package assessment

import (
	"time"

	"github.com/google/uuid"
)

type Assessment struct {
	ID                   uuid.UUID
	PlaceID              uuid.UUID
	AssessedAt           time.Time
	PeriodStart          time.Time
	PeriodEnd            time.Time
	UrgencyScore         int // 0–100, backend-only
	ConfidenceScore      int // 0–100, data freshness + source reliability
	ApplicableIndicators []uuid.UUID
	AffectedGroups       []string
	AssessmentSummary    *string
	CreatedAt            time.Time
}

type AssessmentRequest struct {
	PlaceID              uuid.UUID
	PeriodStart          time.Time
	PeriodEnd            time.Time
	ApplicableIndicators []IndicatorSignal
	DataFreshness        DataFreshness
}

type IndicatorSignal struct {
	IndicatorID    uuid.UUID
	Code           string
	Value          float64
	Trend          *string
	RelevanceScore float64
}

type DataFreshness struct {
	OldestObservationAt time.Time
	SourceReliability   map[string]float64 // source_code -> reliability (0–1)
}

type AffectedGroup struct {
	Group         string
	IndicatorCode string
	Severity      int // 0–100
}
