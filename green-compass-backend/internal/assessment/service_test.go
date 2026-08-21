package assessment

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
)

type mockAssessmentRepo struct {
	stored []*Assessment
}

func (m *mockAssessmentRepo) Store(ctx context.Context, a *Assessment) error {
	a.ID = uuid.New()
	a.CreatedAt = time.Now()
	m.stored = append(m.stored, a)
	return nil
}

func (m *mockAssessmentRepo) GetLatest(ctx context.Context, placeID uuid.UUID) (*Assessment, error) {
	return nil, ErrNotFound
}

func (m *mockAssessmentRepo) GetForPeriod(ctx context.Context, placeID uuid.UUID, periodStart, periodEnd time.Time) (*Assessment, error) {
	return nil, ErrNotFound
}

func newMockAssessmentRepo() *mockAssessmentRepo {
	return &mockAssessmentRepo{stored: make([]*Assessment, 0)}
}

func TestComputeUrgency(t *testing.T) {
	tests := []struct {
		name            string
		indicators      []IndicatorSignal
		wantUrgencyMin  int
		wantUrgencyMax  int
	}{
		{
			name: "high rainfall urgency",
			indicators: []IndicatorSignal{
				{
					IndicatorID:    uuid.New(),
					Code:           "rain_intensity_trend",
					Value:          18.0,
					RelevanceScore: 1.0,
				},
			},
			wantUrgencyMin: 40,
			wantUrgencyMax: 100,
		},
		{
			name: "low soil moisture urgency",
			indicators: []IndicatorSignal{
				{
					IndicatorID:    uuid.New(),
					Code:           "soil_moisture_level",
					Value:          15.0, // 15% = low
					RelevanceScore: 1.0,
				},
			},
			wantUrgencyMin: 60,
			wantUrgencyMax: 100,
		},
		{
			name: "high drought stress urgency",
			indicators: []IndicatorSignal{
				{
					IndicatorID:    uuid.New(),
					Code:           "drought_stress",
					Value:          85.0,
					RelevanceScore: 1.0,
				},
			},
			wantUrgencyMin: 70,
			wantUrgencyMax: 100,
		},
		{
			name: "no indicators",
			indicators: []IndicatorSignal{},
			wantUrgencyMin: 0,
			wantUrgencyMax: 0,
		},
		{
			name: "multiple indicators combined",
			indicators: []IndicatorSignal{
				{
					IndicatorID:    uuid.New(),
					Code:           "rain_intensity_trend",
					Value:          10.0,
					RelevanceScore: 0.8,
				},
				{
					IndicatorID:    uuid.New(),
					Code:           "soil_moisture_level",
					Value:          30.0,
					RelevanceScore: 0.9,
				},
			},
			wantUrgencyMin: 20,
			wantUrgencyMax: 80,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewService(newMockAssessmentRepo(), slog.Default())
			urgency := svc.computeUrgency(tt.indicators)

			if urgency < tt.wantUrgencyMin || urgency > tt.wantUrgencyMax {
				t.Errorf("Urgency: got %d, want in range [%d, %d]", urgency, tt.wantUrgencyMin, tt.wantUrgencyMax)
			}
		})
	}
}

func TestComputeConfidence(t *testing.T) {
	tests := []struct {
		name                 string
		oldestObservationAge time.Duration
		sourceReliability    map[string]float64
		wantConfidenceMin    int
		wantConfidenceMax    int
	}{
		{
			name:                 "fresh data, high reliability",
			oldestObservationAge: 2 * 24 * time.Hour,
			sourceReliability:    map[string]float64{"open_meteo": 0.95, "nasa_power": 0.90},
			wantConfidenceMin:    80,
			wantConfidenceMax:    100,
		},
		{
			name:                 "week-old data, medium reliability",
			oldestObservationAge: 7 * 24 * time.Hour,
			sourceReliability:    map[string]float64{"source": 0.70},
			wantConfidenceMin:    60,
			wantConfidenceMax:    90,
		},
		{
			name:                 "month-old data, low reliability",
			oldestObservationAge: 35 * 24 * time.Hour,
			sourceReliability:    map[string]float64{"source": 0.40},
			wantConfidenceMin:    30,
			wantConfidenceMax:    60,
		},
		{
			name:                 "no source reliability data",
			oldestObservationAge: 3 * 24 * time.Hour,
			sourceReliability:    map[string]float64{},
			wantConfidenceMin:    50,
			wantConfidenceMax:    80,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewService(newMockAssessmentRepo(), slog.Default())
			confidence := svc.computeConfidence(DataFreshness{
				OldestObservationAt: time.Now().Add(-tt.oldestObservationAge),
				SourceReliability:   tt.sourceReliability,
			})

			if confidence < tt.wantConfidenceMin || confidence > tt.wantConfidenceMax {
				t.Errorf("Confidence: got %d, want in range [%d, %d]", confidence, tt.wantConfidenceMin, tt.wantConfidenceMax)
			}
		})
	}
}

func TestDetermineAffectedGroups(t *testing.T) {
	tests := []struct {
		name           string
		indicators     []IndicatorSignal
		wantGroupCount int
		wantHasGroups  []string
	}{
		{
			name: "rain and drought affect farmers",
			indicators: []IndicatorSignal{
				{IndicatorID: uuid.New(), Code: "rain_intensity_trend", Value: 5.0, RelevanceScore: 1.0},
				{IndicatorID: uuid.New(), Code: "drought_stress", Value: 60.0, RelevanceScore: 1.0},
			},
			wantGroupCount: 2, // farmers, water_users
			wantHasGroups:  []string{"farmers", "water_users"},
		},
		{
			name: "water quality affects multiple groups",
			indicators: []IndicatorSignal{
				{IndicatorID: uuid.New(), Code: "water_quality_concern", Value: 1.0, RelevanceScore: 1.0},
			},
			wantGroupCount: 2, // water_users, health_workers
			wantHasGroups:  []string{"water_users", "health_workers"},
		},
		{
			name: "air quality affects vulnerable groups",
			indicators: []IndicatorSignal{
				{IndicatorID: uuid.New(), Code: "air_quality_index", Value: 150.0, RelevanceScore: 1.0},
			},
			wantGroupCount: 3, // health_workers, elderly, children
			wantHasGroups:  []string{"health_workers", "elderly", "children"},
		},
		{
			name:           "no indicators",
			indicators:     []IndicatorSignal{},
			wantGroupCount: 0,
			wantHasGroups:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewService(newMockAssessmentRepo(), slog.Default())
			groups := svc.determineAffectedGroups(tt.indicators)

			if len(groups) != tt.wantGroupCount {
				t.Errorf("Groups count: got %d, want %d", len(groups), tt.wantGroupCount)
			}

			for _, wantGroup := range tt.wantHasGroups {
				found := false
				for _, group := range groups {
					if group == wantGroup {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Missing expected group: %s", wantGroup)
				}
			}
		})
	}
}

func TestAssessIntegration(t *testing.T) {
	repo := newMockAssessmentRepo()
	svc := NewService(repo, slog.Default())

	placeID := uuid.New()
	now := time.Now().UTC()

	assessment, err := svc.Assess(context.Background(), AssessmentRequest{
		PlaceID:     placeID,
		PeriodStart: now.Add(-7 * 24 * time.Hour),
		PeriodEnd:   now,
		ApplicableIndicators: []IndicatorSignal{
			{
				IndicatorID:    uuid.New(),
				Code:           "rain_intensity_trend",
				Value:          8.0,
				RelevanceScore: 0.9,
			},
			{
				IndicatorID:    uuid.New(),
				Code:           "soil_moisture_level",
				Value:          35.0,
				RelevanceScore: 0.8,
			},
		},
		DataFreshness: DataFreshness{
			OldestObservationAt: now.Add(-3 * 24 * time.Hour),
			SourceReliability: map[string]float64{
				"open_meteo": 0.92,
			},
		},
	})

	if err != nil {
		t.Fatalf("Assess failed: %v", err)
	}

	if assessment.PlaceID != placeID {
		t.Errorf("PlaceID: got %v, want %v", assessment.PlaceID, placeID)
	}

	if assessment.UrgencyScore < 0 || assessment.UrgencyScore > 100 {
		t.Errorf("UrgencyScore out of range: %d", assessment.UrgencyScore)
	}

	if assessment.ConfidenceScore < 0 || assessment.ConfidenceScore > 100 {
		t.Errorf("ConfidenceScore out of range: %d", assessment.ConfidenceScore)
	}

	if len(assessment.ApplicableIndicators) != 2 {
		t.Errorf("ApplicableIndicators count: got %d, want 2", len(assessment.ApplicableIndicators))
	}

	if len(assessment.AffectedGroups) == 0 {
		t.Error("AffectedGroups should not be empty")
	}

	if assessment.AssessmentSummary == nil || *assessment.AssessmentSummary == "" {
		t.Error("AssessmentSummary should not be empty")
	}

	if len(repo.stored) != 1 {
		t.Errorf("Repository stored count: got %d, want 1", len(repo.stored))
	}
}
