package updates

import (
	"time"

	"github.com/google/uuid"
)

// Update represents a generated content item ready for user consumption
type Update struct {
	ID           uuid.UUID   `json:"id"`
	PlaceID      uuid.UUID   `json:"place_id"`
	PlaceName    string      `json:"place_name"`
	ContentType  string      `json:"content_type"` // today, forecast, alert
	Headline     string      `json:"headline"`
	BodyText     string      `json:"body_text"`
	CallToAction *string     `json:"call_to_action,omitempty"`
	Indicators   []Indicator `json:"indicators,omitempty"`
	GeneratedAt  time.Time   `json:"generated_at"`
	PeriodStart  time.Time   `json:"period_start"`
	PeriodEnd    time.Time   `json:"period_end"`
}

// Indicator represents a single indicator in an update
type Indicator struct {
	ID             uuid.UUID `json:"id"`
	Code           string    `json:"code"`
	DisplayName    string    `json:"display_name"`
	Value          float64   `json:"value"`
	Unit           string    `json:"unit"`
	Trend          *string   `json:"trend,omitempty"`
	RelevanceScore float64   `json:"relevance_score"`
}

// ListRequest represents query parameters for listing updates
type ListRequest struct {
	UserID  uuid.UUID
	PlaceID *uuid.UUID
	Page    int
	Limit   int
}

// ListResponse is the paginated result
type ListResponse struct {
	Updates []Update `json:"updates"`
	Total   int      `json:"total"`
	Page    int      `json:"page"`
	Limit   int      `json:"limit"`
	HasNext bool     `json:"has_next"`
}

// GetTodayRequest for retrieving today's update
type GetTodayRequest struct {
	UserID uuid.UUID
}

// ExploreRequest represents parameters for exploring indicators
type ExploreRequest struct {
	UserID   uuid.UUID
	Category *string // weather, water, agriculture, air_quality
	PlaceID  *uuid.UUID
	Page     int
	Limit    int
}

// ExploreIndicator represents an indicator in explore view
type ExploreIndicator struct {
	ID              uuid.UUID `json:"id"`
	PlaceID         uuid.UUID `json:"place_id"`
	PlaceName       string    `json:"place_name"`
	Code            string    `json:"code"`
	DisplayName     string    `json:"display_name"`
	Category        string    `json:"category"`
	Value           float64   `json:"value"`
	Unit            string    `json:"unit"`
	Trend           *string   `json:"trend,omitempty"`
	TrendConfidence *float64  `json:"trend_confidence,omitempty"`
	ComputedAt      time.Time `json:"computed_at"`
}

// ExploreResponse is the paginated explore result
type ExploreResponse struct {
	Indicators []ExploreIndicator `json:"indicators"`
	Total      int                `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	HasNext    bool               `json:"has_next"`
}
