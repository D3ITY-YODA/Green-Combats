package applicability

import (
	"context"
	"log/slog"
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

// Filter returns only the applicable indicators for a place
func (s *Service) Filter(ctx context.Context, req FilterRequest) (*FilterResult, error) {
	rules, err := s.repo.ListRulesForPlaceType(ctx, req.PlaceType)
	if err != nil {
		return nil, err
	}

	result := &FilterResult{
		ApplicableIndicators: []ApplicabilityResult{},
		FilteredCount:        0,
	}

	// Build a map of rules for fast lookup
	ruleMap := make(map[string]*Rule)
	for i := range rules {
		ruleMap[rules[i].IndicatorID.String()] = &rules[i]
	}

	for _, ind := range req.Indicators {
		rule, hasRule := ruleMap[ind.IndicatorID.String()]

		// Default: apply indicator if no explicit rule exists
		applicable := true
		relevanceScore := 1.0

		if hasRule {
			applicable = rule.Applicable
			relevanceScore = rule.RelevanceScore
		}

		if !applicable {
			result.FilteredCount++
			continue
		}

		meetsThreshold := true
		if hasRule {
			if rule.MinThreshold != nil && ind.Value < *rule.MinThreshold {
				meetsThreshold = false
			}
			if rule.MaxThreshold != nil && ind.Value > *rule.MaxThreshold {
				meetsThreshold = false
			}
		}

		result.ApplicableIndicators = append(result.ApplicableIndicators, ApplicabilityResult{
			IndicatorID:    ind.IndicatorID,
			IndicatorCode:  ind.Code,
			DisplayName:    ind.DisplayName,
			Applicable:     applicable,
			RelevanceScore: relevanceScore,
			MeetsThreshold: meetsThreshold,
			Value:          ind.Value,
			Rule:           rule,
		})
	}

	return result, nil
}
