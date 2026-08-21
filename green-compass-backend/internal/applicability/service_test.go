package applicability

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
)

type mockAppRepo struct {
	rules map[string]*Rule
}

func (m *mockAppRepo) GetRuleForIndicatorAndPlaceType(ctx context.Context, indicatorID uuid.UUID, placeType string) (*Rule, error) {
	key := indicatorID.String() + ":" + placeType
	if rule, ok := m.rules[key]; ok {
		return rule, nil
	}
	return nil, ErrRuleNotFound
}

func (m *mockAppRepo) ListRulesForPlaceType(ctx context.Context, placeType string) ([]Rule, error) {
	var rules []Rule
	for _, rule := range m.rules {
		if rule.PlaceType == placeType {
			rules = append(rules, *rule)
		}
	}
	return rules, nil
}

func newMockAppRepo() *mockAppRepo {
	return &mockAppRepo{
		rules: make(map[string]*Rule),
	}
}

func (m *mockAppRepo) addRule(rule *Rule) {
	key := rule.IndicatorID.String() + ":" + rule.PlaceType
	m.rules[key] = rule
}

func TestFilterApplicableIndicators(t *testing.T) {
	indID1 := uuid.New()
	indID2 := uuid.New()
	indID3 := uuid.New()

	tests := []struct {
		name         string
		placeType    string
		indicators   []IndicatorData
		rules        []*Rule
		wantCount    int
		wantFiltered int
	}{
		{
			name:      "all applicable for agricultural place",
			placeType: "agricultural",
			indicators: []IndicatorData{
				{IndicatorID: indID1, Code: "rain_intensity_trend", DisplayName: "Rain Intensity", Value: 5.0},
				{IndicatorID: indID2, Code: "soil_moisture_level", DisplayName: "Soil Moisture", Value: 40.0},
			},
			rules: []*Rule{
				{IndicatorID: indID1, PlaceType: "agricultural", Applicable: true, RelevanceScore: 1.0},
				{IndicatorID: indID2, PlaceType: "agricultural", Applicable: true, RelevanceScore: 1.0},
			},
			wantCount:    2,
			wantFiltered: 0,
		},
		{
			name:      "some filtered out",
			placeType: "urban",
			indicators: []IndicatorData{
				{IndicatorID: indID1, Code: "rain_intensity_trend", DisplayName: "Rain Intensity", Value: 5.0},
				{IndicatorID: indID2, Code: "soil_moisture_level", DisplayName: "Soil Moisture", Value: 40.0},
				{IndicatorID: indID3, Code: "drought_stress", DisplayName: "Drought", Value: 25.0},
			},
			rules: []*Rule{
				{IndicatorID: indID1, PlaceType: "urban", Applicable: true, RelevanceScore: 0.8},
				{IndicatorID: indID2, PlaceType: "urban", Applicable: false, RelevanceScore: 0.3},
				{IndicatorID: indID3, PlaceType: "urban", Applicable: true, RelevanceScore: 0.5},
			},
			wantCount:    2,
			wantFiltered: 1,
		},
		{
			name:      "threshold filtering",
			placeType: "water_community",
			indicators: []IndicatorData{
				{IndicatorID: indID1, Code: "rain_intensity_trend", DisplayName: "Rain Intensity", Value: 2.0},
				{IndicatorID: indID2, Code: "rain_frequency", DisplayName: "Rain Frequency", Value: 8.0},
			},
			rules: []*Rule{
				{IndicatorID: indID1, PlaceType: "water_community", Applicable: true, RelevanceScore: 1.0, MinThreshold: floatPtr(3.0)},
				{IndicatorID: indID2, PlaceType: "water_community", Applicable: true, RelevanceScore: 1.0, MaxThreshold: floatPtr(5.0)},
			},
			wantCount:    2,
			wantFiltered: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockAppRepo()
			for _, rule := range tt.rules {
				repo.addRule(rule)
			}
			svc := NewService(repo, slog.Default())

			result, err := svc.Filter(context.Background(), FilterRequest{
				PlaceID:    uuid.New(),
				PlaceType:  tt.placeType,
				Indicators: tt.indicators,
			})

			if err != nil {
				t.Fatalf("Filter failed: %v", err)
			}

			if len(result.ApplicableIndicators) != tt.wantCount {
				t.Errorf("ApplicableIndicators count: got %d, want %d", len(result.ApplicableIndicators), tt.wantCount)
			}

			if result.FilteredCount != tt.wantFiltered {
				t.Errorf("FilteredCount: got %d, want %d", result.FilteredCount, tt.wantFiltered)
			}
		})
	}
}

func TestThresholdChecking(t *testing.T) {
	indID := uuid.New()

	tests := []struct {
		name               string
		value              float64
		minThreshold       *float64
		maxThreshold       *float64
		wantMeetsThreshold bool
	}{
		{
			name:               "value above min threshold",
			value:              5.0,
			minThreshold:       floatPtr(3.0),
			maxThreshold:       nil,
			wantMeetsThreshold: true,
		},
		{
			name:               "value below min threshold",
			value:              2.0,
			minThreshold:       floatPtr(3.0),
			maxThreshold:       nil,
			wantMeetsThreshold: false,
		},
		{
			name:               "value below max threshold",
			value:              4.0,
			minThreshold:       nil,
			maxThreshold:       floatPtr(5.0),
			wantMeetsThreshold: true,
		},
		{
			name:               "value above max threshold",
			value:              6.0,
			minThreshold:       nil,
			maxThreshold:       floatPtr(5.0),
			wantMeetsThreshold: false,
		},
		{
			name:               "value within range",
			value:              4.5,
			minThreshold:       floatPtr(3.0),
			maxThreshold:       floatPtr(5.0),
			wantMeetsThreshold: true,
		},
		{
			name:               "value outside range (too low)",
			value:              2.5,
			minThreshold:       floatPtr(3.0),
			maxThreshold:       floatPtr(5.0),
			wantMeetsThreshold: false,
		},
		{
			name:               "value outside range (too high)",
			value:              5.5,
			minThreshold:       floatPtr(3.0),
			maxThreshold:       floatPtr(5.0),
			wantMeetsThreshold: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockAppRepo()
			rule := &Rule{
				IndicatorID:    indID,
				PlaceType:      "test",
				Applicable:     true,
				RelevanceScore: 1.0,
				MinThreshold:   tt.minThreshold,
				MaxThreshold:   tt.maxThreshold,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			repo.addRule(rule)

			svc := NewService(repo, slog.Default())
			result, err := svc.Filter(context.Background(), FilterRequest{
				PlaceID:   uuid.New(),
				PlaceType: "test",
				Indicators: []IndicatorData{
					{IndicatorID: indID, Code: "test_ind", DisplayName: "Test", Value: tt.value},
				},
			})

			if err != nil {
				t.Fatalf("Filter failed: %v", err)
			}

			if len(result.ApplicableIndicators) == 0 {
				t.Fatal("No applicable indicators returned")
			}

			if result.ApplicableIndicators[0].MeetsThreshold != tt.wantMeetsThreshold {
				t.Errorf("MeetsThreshold: got %v, want %v", result.ApplicableIndicators[0].MeetsThreshold, tt.wantMeetsThreshold)
			}
		})
	}
}

// Helpers
func floatPtr(v float64) *float64 {
	return &v
}
