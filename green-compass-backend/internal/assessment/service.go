package assessment

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type assessmentRepo interface {
	Store(ctx context.Context, a *Assessment) error
	GetLatest(ctx context.Context, placeID uuid.UUID) (*Assessment, error)
	GetForPeriod(ctx context.Context, placeID uuid.UUID, periodStart, periodEnd time.Time) (*Assessment, error)
}

type Service struct {
	repo   assessmentRepo
	logger *slog.Logger
}

func NewService(repo assessmentRepo, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// Assess computes backend-only scoring for a place based on applicable indicators
func (s *Service) Assess(ctx context.Context, req AssessmentRequest) (*Assessment, error) {
	urgencyScore := s.computeUrgency(req.ApplicableIndicators)
	confidenceScore := s.computeConfidence(req.DataFreshness)
	affectedGroups := s.determineAffectedGroups(req.ApplicableIndicators)

	assessment := &Assessment{
		ID:                   uuid.New(),
		PlaceID:              req.PlaceID,
		AssessedAt:           time.Now().UTC(),
		PeriodStart:          req.PeriodStart,
		PeriodEnd:            req.PeriodEnd,
		UrgencyScore:         urgencyScore,
		ConfidenceScore:      confidenceScore,
		ApplicableIndicators: toUUIDArray(req.ApplicableIndicators),
		AffectedGroups:       affectedGroups,
		AssessmentSummary:    s.generateSummary(req.ApplicableIndicators, urgencyScore),
	}

	if err := s.repo.Store(ctx, assessment); err != nil {
		return nil, err
	}

	return assessment, nil
}

// computeUrgency scores 0–100 based on indicator values and trends
func (s *Service) computeUrgency(indicators []IndicatorSignal) int {
	if len(indicators) == 0 {
		return 0
	}

	totalScore := 0.0
	for _, ind := range indicators {
		score := s.indicatorUrgency(ind.Code, ind.Value, ind.Trend, ind.RelevanceScore)
		totalScore += score * ind.RelevanceScore
	}

	avgScore := totalScore / float64(len(indicators))
	if avgScore > 100 {
		avgScore = 100
	}
	return int(avgScore)
}

// indicatorUrgency returns a 0–100 score for a single indicator
func (s *Service) indicatorUrgency(code string, value float64, trend *string, relScore float64) float64 {
	switch code {
	case "rain_intensity_trend":
		// Higher rainfall = higher urgency for flood risk
		return s.normalizeToScore(value, 0, 20) * 0.7
	case "rain_frequency":
		// More rainy days = higher urgency
		return s.normalizeToScore(value, 0, 10) * 0.6
	case "soil_moisture_level":
		// Low moisture = high urgency (drought)
		return (100.0 - s.normalizeToScore(value, 0, 100)) * 0.8
	case "drought_stress":
		// Drought index 0–100 maps directly to urgency
		return value * 0.9
	case "temperature_trend":
		// Extreme temps increase urgency
		if value > 35 || value < 0 {
			return 70
		}
		return 30
	case "wind_speed_trend":
		// High winds = moderate urgency
		return s.normalizeToScore(value, 0, 20) * 0.5
	default:
		return 20 // Default low urgency for unknown indicators
	}
}

// normalizeToScore converts a value to 0–100 scale
func (s *Service) normalizeToScore(value, min, max float64) float64 {
	if max <= min {
		return 0
	}
	score := (value - min) / (max - min) * 100
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}
	return score
}

// computeConfidence scores 0–100 based on data freshness and source reliability
func (s *Service) computeConfidence(freshness DataFreshness) int {
	// Freshness component: data age
	ageDays := time.Since(freshness.OldestObservationAt).Hours() / 24
	freshnessScore := 100.0
	if ageDays > 30 {
		freshnessScore = 40
	} else if ageDays > 7 {
		freshnessScore = 70
	}

	// Reliability component: average source reliability
	reliabilityScore := 0.0
	if len(freshness.SourceReliability) > 0 {
		sum := 0.0
		for _, rel := range freshness.SourceReliability {
			sum += rel
		}
		reliabilityScore = (sum / float64(len(freshness.SourceReliability))) * 100
	} else {
		reliabilityScore = 50 // Neutral default
	}

	// Weighted average: 60% freshness, 40% reliability
	confidence := (freshnessScore * 0.6) + (reliabilityScore * 0.4)
	if confidence > 100 {
		confidence = 100
	}
	return int(confidence)
}

// determineAffectedGroups returns list of groups affected by indicator conditions
func (s *Service) determineAffectedGroups(indicators []IndicatorSignal) []string {
	groups := make(map[string]bool)

	for _, ind := range indicators {
		switch ind.Code {
		case "rain_intensity_trend", "rain_frequency":
			groups["farmers"] = true
		case "soil_moisture_level", "drought_stress":
			groups["farmers"] = true
			groups["water_users"] = true
		case "water_quality_concern":
			groups["water_users"] = true
			groups["health_workers"] = true
		case "air_quality_index":
			groups["health_workers"] = true
			groups["elderly"] = true
			groups["children"] = true
		}
	}

	// Convert map to sorted slice
	result := make([]string, 0, len(groups))
	for group := range groups {
		result = append(result, group)
	}
	return result
}

// generateSummary creates a brief summary (backend use, not exposed in API)
func (s *Service) generateSummary(indicators []IndicatorSignal, urgency int) *string {
	if len(indicators) == 0 {
		return nil
	}

	urgencyLabel := "low"
	if urgency > 70 {
		urgencyLabel = "high"
	} else if urgency > 40 {
		urgencyLabel = "moderate"
	}

	summary := fmt.Sprintf("Assessment: %s urgency based on %d applicable indicator(s)", urgencyLabel, len(indicators))
	return &summary
}

// toUUIDArray converts IndicatorSignal slice to UUID array
func toUUIDArray(signals []IndicatorSignal) []uuid.UUID {
	arr := make([]uuid.UUID, len(signals))
	for i, sig := range signals {
		arr[i] = sig.IndicatorID
	}
	return arr
}
