package content

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Service struct {
	repo   *Repository
	logger *slog.Logger
}

func NewService(repo *Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// Generate produces plain-language content from assessment + indicators
func (s *Service) Generate(ctx context.Context, req GenerateRequest) (*Content, error) {
	headline := s.generateHeadline(req)
	bodyText := s.generateBody(req)
	callToAction := s.generateCallToAction(req)

	sourceIndicatorIDs := make(pq.UUIDArray, 0, len(req.ApplicableIndicators))
	for _, ind := range req.ApplicableIndicators {
		sourceIndicatorIDs = append(sourceIndicatorIDs, ind.IndicatorID)
	}

	content := &Content{
		ID:               uuid.New(),
		PlaceID:          req.PlaceID,
		GeneratedAt:      time.Now().UTC(),
		PeriodStart:      req.PeriodStart,
		PeriodEnd:        req.PeriodEnd,
		ContentType:      req.ContentType,
		Language:         req.Language,
		Headline:         headline,
		BodyText:         bodyText,
		CallToAction:     callToAction,
		SourceIndicators: sourceIndicatorIDs,
	}

	if err := s.repo.Store(ctx, content); err != nil {
		return nil, err
	}

	return content, nil
}

// generateHeadline creates the headline based on urgency and content type
func (s *Service) generateHeadline(req GenerateRequest) string {
	switch req.ContentType {
	case "alert":
		return s.generateAlertHeadline(req)
	case "forecast":
		return s.generateForecastHeadline(req)
	default: // "today"
		return s.generateTodayHeadline(req)
	}
}

func (s *Service) generateAlertHeadline(req GenerateRequest) string {
	if req.Assessment.UrgencyScore > 70 {
		return fmt.Sprintf("🚨 High Alert: Critical conditions in %s", req.PlaceName)
	} else if req.Assessment.UrgencyScore > 40 {
		return fmt.Sprintf("⚠️  Important Update: Conditions changing in %s", req.PlaceName)
	}
	return fmt.Sprintf("ℹ️ Information: Local conditions update for %s", req.PlaceName)
}

func (s *Service) generateForecastHeadline(req GenerateRequest) string {
	if len(req.ApplicableIndicators) == 0 {
		return fmt.Sprintf("Forecast for %s", req.PlaceName)
	}
	topIndicator := req.ApplicableIndicators[0]
	return fmt.Sprintf("Forecast: %s expected in %s", topIndicator.DisplayName, req.PlaceName)
}

func (s *Service) generateTodayHeadline(req GenerateRequest) string {
	if len(req.ApplicableIndicators) == 0 {
		return fmt.Sprintf("Today in %s", req.PlaceName)
	}
	topIndicator := req.ApplicableIndicators[0]
	return fmt.Sprintf("Today: %s in %s", topIndicator.DisplayName, req.PlaceName)
}

// generateBody creates the main content text
func (s *Service) generateBody(req GenerateRequest) string {
	var parts []string

	// Opening based on urgency
	if req.Assessment.UrgencyScore > 70 {
		parts = append(parts, fmt.Sprintf("Conditions in %s require immediate attention.", req.PlaceName))
	} else if req.Assessment.UrgencyScore > 40 {
		parts = append(parts, fmt.Sprintf("Changes are occurring in %s that may affect local activities.", req.PlaceName))
	} else {
		parts = append(parts, fmt.Sprintf("The following conditions are noted in %s:", req.PlaceName))
	}

	// Indicator details
	if len(req.ApplicableIndicators) > 0 {
		parts = append(parts, "\nObservations:")
		for _, ind := range req.ApplicableIndicators {
			desc := s.describeIndicator(ind)
			if desc != "" {
				parts = append(parts, fmt.Sprintf("• %s", desc))
			}
		}
	}

	// Affected groups
	if len(req.Assessment.AffectedGroups) > 0 {
		parts = append(parts, fmt.Sprintf("\nThis may affect: %s.", strings.Join(req.Assessment.AffectedGroups, ", ")))
	}

	// Confidence note (internal-only for now, could be visible)
	if req.Assessment.ConfidenceScore < 50 {
		parts = append(parts, "\nNote: Information based on limited recent data.")
	}

	return strings.Join(parts, "\n")
}

// describeIndicator produces a conversational description of an indicator
func (s *Service) describeIndicator(ind IndicatorData) string {
	switch ind.Code {
	case "rain_intensity_trend":
		if ind.Value > 10 {
			return fmt.Sprintf("Heavy rainfall at %.1f%s", ind.Value, ind.Unit)
		} else if ind.Value > 5 {
			return fmt.Sprintf("Moderate rainfall at %.1f%s", ind.Value, ind.Unit)
		}
		return fmt.Sprintf("Light rainfall at %.1f%s", ind.Value, ind.Unit)

	case "rain_frequency":
		days := int(ind.Value)
		if days > 5 {
			return fmt.Sprintf("Frequent rain with %d rainy days", days)
		}
		return fmt.Sprintf("Rain on approximately %d days", days)

	case "soil_moisture_level":
		if ind.Value < 30 {
			return fmt.Sprintf("Soil is very dry at %.1f%% moisture", ind.Value)
		} else if ind.Value < 50 {
			return fmt.Sprintf("Soil is moderately dry at %.1f%% moisture", ind.Value)
		}
		return fmt.Sprintf("Soil moisture adequate at %.1f%%", ind.Value)

	case "drought_stress":
		if ind.Value > 70 {
			return fmt.Sprintf("Severe drought stress indicated (index: %.0f)", ind.Value)
		} else if ind.Value > 40 {
			return fmt.Sprintf("Moderate drought conditions developing (index: %.0f)", ind.Value)
		}
		return fmt.Sprintf("Low drought stress (index: %.0f)", ind.Value)

	case "temperature_trend":
		if ind.Value > 35 {
			return fmt.Sprintf("Very high temperatures at %.1f°%s", ind.Value, ind.Unit)
		} else if ind.Value < 5 {
			return fmt.Sprintf("Very cold conditions at %.1f°%s", ind.Value, ind.Unit)
		}
		return fmt.Sprintf("Temperature at %.1f°%s", ind.Value, ind.Unit)

	case "wind_speed_trend":
		if ind.Value > 15 {
			return fmt.Sprintf("Strong winds at %.1f %s", ind.Value, ind.Unit)
		}
		return fmt.Sprintf("Wind speed %.1f %s", ind.Value, ind.Unit)

	default:
		return fmt.Sprintf("%s: %.1f %s", ind.DisplayName, ind.Value, ind.Unit)
	}
}

// generateCallToAction produces actionable guidance (if needed)
func (s *Service) generateCallToAction(req GenerateRequest) *string {
	if req.Assessment.UrgencyScore <= 30 {
		return nil // No action needed for low urgency
	}

	var action string
	switch req.ContentType {
	case "alert":
		if req.Assessment.UrgencyScore > 70 {
			action = "Take immediate protective measures. Contact local authorities if needed."
		} else {
			action = "Monitor conditions and prepare contingency plans."
		}

	case "forecast":
		action = "Plan activities accordingly based on expected conditions."

	default:
		action = "Stay informed of any changes in local conditions."
	}

	return &action
}
