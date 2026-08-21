package broker

import (
	"time"

	"github.com/google/uuid"
)

// Topic constants
const (
	TopicContentGenerated  = "content.generated"
	TopicObservationSubmit = "observation.submitted"
	TopicReportVerified    = "report.verified"
)

// ContentGeneratedEvent is published when new plain-language content is produced for a place.
type ContentGeneratedEvent struct {
	ContentID uuid.UUID
	PlaceID   uuid.UUID
	PlaceName string
	Language  string
	CreatedAt time.Time
}

func (e ContentGeneratedEvent) Topic() string { return TopicContentGenerated }

// ObservationSubmittedEvent is published when a community user submits an observation.
type ObservationSubmittedEvent struct {
	ObservationID uuid.UUID
	UserID        uuid.UUID
	PlaceID       uuid.UUID
	SubmittedAt   time.Time
}

func (e ObservationSubmittedEvent) Topic() string { return TopicObservationSubmit }

// ReportVerifiedEvent is published when an institutional user verifies a report.
type ReportVerifiedEvent struct {
	ObservationID uuid.UUID
	OrgID         uuid.UUID
	VerifiedBy    uuid.UUID
	VerifiedAt    time.Time
}

func (e ReportVerifiedEvent) Topic() string { return TopicReportVerified }
