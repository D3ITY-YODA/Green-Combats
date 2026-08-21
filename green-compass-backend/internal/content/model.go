package content

import (
	"time"

	"github.com/google/uuid"
)

type Content struct {
	ID               uuid.UUID
	PlaceID          uuid.UUID
	GeneratedAt      time.Time
	PeriodStart      time.Time
	PeriodEnd        time.Time
	ContentType      string // today, forecast, alert
	Language         string // e.g. "en", "sw"
	Headline         string
	BodyText         string
	CallToAction     *string
	SourceIndicators []uuid.UUID
	CreatedAt        time.Time
}

type GenerateRequest struct {
	PlaceID              uuid.UUID
	PlaceName            string
	PlaceType            string
	PeriodStart          time.Time
	PeriodEnd            time.Time
	ContentType          string // today, forecast, alert
	Language             string
	Assessment           AssessmentData
	ApplicableIndicators []IndicatorData
}

type AssessmentData struct {
	UrgencyScore    int
	ConfidenceScore int
	AffectedGroups  []string
	Summary         string
}

type IndicatorData struct {
	IndicatorID    uuid.UUID
	Code           string
	DisplayName    string
	Value          float64
	Unit           string
	Trend          *string
	RelevanceScore float64
}

type GenerateResult struct {
	Content *Content
	Error   error
}
