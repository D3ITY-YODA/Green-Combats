package content

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
)

type mockContentRepo struct {
	stored []*Content
}

func (m *mockContentRepo) Store(ctx context.Context, c *Content) error {
	c.ID = uuid.New()
	c.CreatedAt = time.Now()
	m.stored = append(m.stored, c)
	return nil
}

func (m *mockContentRepo) GetLatestByType(ctx context.Context, placeID uuid.UUID, contentType string) (*Content, error) {
	return nil, ErrNotFound
}

func (m *mockContentRepo) ListForPlace(ctx context.Context, placeID uuid.UUID, limit int) ([]Content, error) {
	return nil, ErrNotFound
}

func newMockContentRepo() *mockContentRepo {
	return &mockContentRepo{stored: make([]*Content, 0)}
}

func TestGenerateHeadline(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		urgency     int
		placeName   string
		wantPrefix  string
	}{
		{
			name:        "high alert headline",
			contentType: "alert",
			urgency:     80,
			placeName:   "Kampala",
			wantPrefix:  "🚨 High Alert",
		},
		{
			name:        "moderate alert headline",
			contentType: "alert",
			urgency:     50,
			placeName:   "Kampala",
			wantPrefix:  "⚠️  Important Update",
		},
		{
			name:        "low urgency alert headline",
			contentType: "alert",
			urgency:     20,
			placeName:   "Kampala",
			wantPrefix:  "ℹ️ Information",
		},
		{
			name:        "forecast headline",
			contentType: "forecast",
			urgency:     50,
			placeName:   "Mwanza",
			wantPrefix:  "Forecast",
		},
		{
			name:        "today headline",
			contentType: "today",
			urgency:     50,
			placeName:   "Dar es Salaam",
			wantPrefix:  "Today",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewService(newMockContentRepo(), slog.Default())
			headline := svc.generateHeadline(GenerateRequest{
				PlaceName:   tt.placeName,
				ContentType: tt.contentType,
				Assessment: AssessmentData{
					UrgencyScore: tt.urgency,
				},
			})

			if !matchesPrefix(headline, tt.wantPrefix) {
				t.Errorf("Headline: got %q, want prefix %q", headline, tt.wantPrefix)
			}
		})
	}
}

func TestDescribeIndicator(t *testing.T) {
	tests := []struct {
		name               string
		code               string
		value              float64
		wantDescSubstring  string
	}{
		{
			name:              "heavy rainfall",
			code:              "rain_intensity_trend",
			value:             15.0,
			wantDescSubstring: "Heavy rainfall",
		},
		{
			name:              "very dry soil",
			code:              "soil_moisture_level",
			value:             20.0,
			wantDescSubstring: "very dry",
		},
		{
			name:              "severe drought",
			code:              "drought_stress",
			value:             80.0,
			wantDescSubstring: "Severe drought",
		},
		{
			name:              "very high temperature",
			code:              "temperature_trend",
			value:             37.5,
			wantDescSubstring: "Very high temperatures",
		},
		{
			name:              "strong winds",
			code:              "wind_speed_trend",
			value:             18.0,
			wantDescSubstring: "Strong winds",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewService(newMockContentRepo(), slog.Default())
			desc := svc.describeIndicator(IndicatorData{
				Code:        tt.code,
				Value:       tt.value,
				DisplayName: "Test Indicator",
				Unit:        "units",
			})

			if !contains(desc, tt.wantDescSubstring) {
				t.Errorf("Description: got %q, want substring %q", desc, tt.wantDescSubstring)
			}
		})
	}
}

func TestGenerateCallToAction(t *testing.T) {
	tests := []struct {
		name           string
		contentType    string
		urgency        int
		wantNil        bool
		wantSubstring  string
	}{
		{
			name:          "low urgency no action",
			contentType:   "today",
			urgency:       20,
			wantNil:       true,
		},
		{
			name:          "high alert action",
			contentType:   "alert",
			urgency:       80,
			wantNil:       false,
			wantSubstring: "immediate protective measures",
		},
		{
			name:          "moderate alert action",
			contentType:   "alert",
			urgency:       50,
			wantNil:       false,
			wantSubstring: "Monitor conditions",
		},
		{
			name:          "forecast action",
			contentType:   "forecast",
			urgency:       40,
			wantNil:       false,
			wantSubstring: "Plan activities",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewService(newMockContentRepo(), slog.Default())
			action := svc.generateCallToAction(GenerateRequest{
				ContentType: tt.contentType,
				Assessment: AssessmentData{
					UrgencyScore: tt.urgency,
				},
			})

			if (action == nil) != tt.wantNil {
				t.Errorf("Call to action: got nil=%v, want nil=%v", action == nil, tt.wantNil)
			}

			if !tt.wantNil && action != nil {
				if !contains(*action, tt.wantSubstring) {
					t.Errorf("Action: got %q, want substring %q", *action, tt.wantSubstring)
				}
			}
		})
	}
}

func TestGenerateIntegration(t *testing.T) {
	repo := newMockContentRepo()
	svc := NewService(repo, slog.Default())

	placeID := uuid.New()
	now := time.Now().UTC()

	content, err := svc.Generate(context.Background(), GenerateRequest{
		PlaceID:     placeID,
		PlaceName:   "Nairobi",
		PlaceType:   "urban",
		PeriodStart: now.Add(-24 * time.Hour),
		PeriodEnd:   now,
		ContentType: "today",
		Language:    "en",
		Assessment: AssessmentData{
			UrgencyScore:    60,
			ConfidenceScore: 85,
			AffectedGroups:  []string{"farmers", "water_users"},
			Summary:         "Moderate conditions",
		},
		ApplicableIndicators: []IndicatorData{
			{
				IndicatorID:    uuid.New(),
				Code:           "rain_intensity_trend",
				DisplayName:    "Rain Intensity",
				Value:          8.5,
				Unit:           "mm",
				Trend:          strPtr("increasing"),
				RelevanceScore: 0.9,
			},
			{
				IndicatorID:    uuid.New(),
				Code:           "soil_moisture_level",
				DisplayName:    "Soil Moisture",
				Value:          45.0,
				Unit:           "%",
				RelevanceScore: 0.8,
			},
		},
	})

	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if content.PlaceID != placeID {
		t.Errorf("PlaceID: got %v, want %v", content.PlaceID, placeID)
	}

	if content.ContentType != "today" {
		t.Errorf("ContentType: got %s, want today", content.ContentType)
	}

	if content.Language != "en" {
		t.Errorf("Language: got %s, want en", content.Language)
	}

	if content.Headline == "" {
		t.Error("Headline should not be empty")
	}

	if content.BodyText == "" {
		t.Error("BodyText should not be empty")
	}

	if len(content.SourceIndicators) != 2 {
		t.Errorf("SourceIndicators count: got %d, want 2", len(content.SourceIndicators))
	}

	// Verify body text contains indicator info
	if !contains(content.BodyText, "Observations") {
		t.Errorf("Body should mention observations: %s", content.BodyText)
	}

	if !contains(content.BodyText, "farmers") && !contains(content.BodyText, "water_users") {
		t.Errorf("Body should mention affected groups: %s", content.BodyText)
	}

	if len(repo.stored) != 1 {
		t.Errorf("Repository stored count: got %d, want 1", len(repo.stored))
	}
}

// Helpers
func strPtr(s string) *string {
	return &s
}

func contains(str, substr string) bool {
	return len(str) > 0 && len(substr) > 0 && (str == substr || len(str) > len(substr) && (str[:len(substr)] == substr || str[len(str)-len(substr):] == substr || indexAny(str, substr) >= 0))
}

func indexAny(str, substr string) int {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func matchesPrefix(str, prefix string) bool {
	return len(str) >= len(prefix) && str[:len(prefix)] == prefix
}
