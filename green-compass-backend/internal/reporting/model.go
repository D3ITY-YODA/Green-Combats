package reporting

import (
	"time"

	"github.com/google/uuid"
)

// DeliveryStats represents notification delivery metrics for a place.
type DeliveryStats struct {
	PlaceID          uuid.UUID `json:"place_id"`
	TotalSent        int       `json:"total_sent"`
	TotalFailed      int       `json:"total_failed"`
	ByChannel        map[string]int `json:"by_channel"`
	PeriodStart      time.Time `json:"period_start"`
	PeriodEnd        time.Time `json:"period_end"`
}

// InteractionStats represents user interaction metrics.
type InteractionStats struct {
	PlaceID          uuid.UUID `json:"place_id"`
	TotalViews       int       `json:"total_views"`
	TotalClicks      int       `json:"total_clicks"`
	UniqueUsers      int       `json:"unique_users"`
	PeriodStart      time.Time `json:"period_start"`
	PeriodEnd        time.Time `json:"period_end"`
}

// ReportStats represents community report metrics.
type ReportStats struct {
	OrgID            uuid.UUID `json:"org_id"`
	TotalSubmitted   int       `json:"total_submitted"`
	TotalVerified    int       `json:"total_verified"`
	TotalRejected    int       `json:"total_rejected"`
	PendingReview    int       `json:"pending_review"`
	PeriodStart      time.Time `json:"period_start"`
	PeriodEnd        time.Time `json:"period_end"`
}

// DashboardResponse combines all metrics for an institutional dashboard.
type DashboardResponse struct {
	Delivery    *DeliveryStats    `json:"delivery,omitempty"`
	Interaction *InteractionStats `json:"interaction,omitempty"`
	Reports     *ReportStats      `json:"reports,omitempty"`
}

// QueryRequest is the input for metrics queries.
type QueryRequest struct {
	OrgID       *uuid.UUID
	PlaceID     *uuid.UUID
	PeriodStart time.Time
	PeriodEnd   time.Time
}
